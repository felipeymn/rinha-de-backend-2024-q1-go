package types

import (
	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type StatementAccount struct {
	Balance   pgtype.Int4 `json:"total"`
	Timestamp string      `json:"data_extrato"`
	Limit     pgtype.Int4 `json:"limite"`
}

type StatementTransaction struct {
	Amount      int32              `json:"valor"`
	Operation   sqlc.OperationEnum `json:"tipo"`
	Description string             `json:"descricao"`
	Timestamp   pgtype.Timestamp   `json:"realizada_em"`
}

type StatementResponse struct {
	Account      StatementAccount       `json:"saldo"`
	Transactions []StatementTransaction `json:"ultimas_transacoes"`
}
