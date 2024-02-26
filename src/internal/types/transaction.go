package types

import (
	"errors"

	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Transaction struct {
	Amount      int32              `json:"valor"`
	Operation   sqlc.OperationEnum `json:"tipo"`
	Description string             `json:"descricao"`
}

type TransactionResponse struct {
	Limit   pgtype.Int4 `json:"limite"`
	Balance int32       `json:"saldo"`
}

func (t Transaction) IsValid() (err error) {
	if t.Operation != "c" && t.Operation != "d" {
		err = errors.New("Type field must be 'c' (credit) or 'd' (debit)")
	}
	if t.Amount < 0 {
		err = errors.New("Value field must be a positive integer")
	}
	if len(t.Description) < 1 || len(t.Description) > 10 {
		err = errors.New("Description field must have 1 to 10 characters")
	}
	return
}
