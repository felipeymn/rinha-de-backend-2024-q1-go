package database

import (
	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/database/sqlc"
	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

type Storage interface {
	GetStatement(int32) (sqlc.GetAccountRow, []sqlc.GetTransactionsRow, error)
	SaveTransaction(int32, types.Transaction) (sqlc.UpdateAccountRow, error)
}
