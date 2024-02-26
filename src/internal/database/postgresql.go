package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/database/sqlc"
	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/types"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgreSQL() *Store {
	config := createConfig()

	ctx := context.Background()
	pool, err := createConnectionPoolWithRetry(
		ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil
	}
	return &Store{pool, sqlc.New(pool)}
}

func createConfig() *pgxpool.Config {
	config, err := pgxpool.ParseConfig(os.Getenv("RINHA_DB_CONN_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse connection string: %v\n", err)
		return nil
	}
	config.MaxConns, config.MinConns = getConnConfigFromEnv()
	return config
}

func getEnvInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getConnConfigFromEnv() (maxConns int32, minConns int32) {
	maxConns = int32(getEnvInt("MAX_CONNS", 20))
	minConns = int32(getEnvInt("MIN_CONNS", 20))
	return maxConns, minConns
}

func createConnectionPoolWithRetry(
	ctx context.Context,
	config *pgxpool.Config,
) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for attempt := 1; attempt <= 20; attempt++ {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			// Ping to ensure the connection is valid
			if err := pool.Ping(ctx); err == nil {
				fmt.Println("connected to database")
				return pool, nil
			}
		}
		backoff := time.Duration(3) * time.Second
		time.Sleep(backoff)
	}
	return nil, fmt.Errorf("failed to create connection pool after retrying")
}

func (p *Store) GetStatement(
	id int32,
) (account sqlc.GetAccountRow, transactions []sqlc.GetTransactionsRow, err error) {
	ctx := context.Background()
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return account, transactions, err
	}
	defer tx.Rollback(ctx)

	qtx := p.Queries.WithTx(tx)
	account, err = qtx.GetAccount(ctx, id)
	if err != nil {
		return account, transactions, err
	}
	transactions, err = qtx.GetTransactions(ctx, pgtype.Int4{Int32: id, Valid: true})
	if err != nil {
		return account, transactions, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return account, transactions, err
	}

	return account, transactions, nil
}

func (p *Store) SaveTransaction(
	id int32,
	transaction types.Transaction,
) (account sqlc.UpdateAccountRow, err error) {
	ctx := context.Background()
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return account, err
	}
	defer tx.Rollback(ctx)

	qtx := p.Queries.WithTx(tx)
	account, err = qtx.UpdateAccount(
		ctx,
		sqlc.UpdateAccountParams{
			Balance: pgtype.Int4{
				Int32: creditOrDebit(transaction.Operation, transaction.Amount),
				Valid: true,
			},
			ID: id,
		},
	)
	if err != nil {
		return account, err
	}
	err = qtx.CreateTransaction(
		ctx,
		sqlc.CreateTransactionParams{
			AccountID:   pgtype.Int4{Int32: id, Valid: true},
			Amount:      transaction.Amount,
			Operation:   sqlc.NullOperationEnum{OperationEnum: transaction.Operation, Valid: true},
			Description: transaction.Description,
		},
	)
	if err != nil {
		return account, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return account, err
	}

	return account, nil
}

func creditOrDebit(operation sqlc.OperationEnum, amount int32) int32 {
	if operation == sqlc.OperationEnumD {
		return -(amount)
	}
	return amount
}
