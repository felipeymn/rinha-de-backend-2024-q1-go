// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OperationEnum string

const (
	OperationEnumC OperationEnum = "c"
	OperationEnumD OperationEnum = "d"
)

func (e *OperationEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OperationEnum(s)
	case string:
		*e = OperationEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for OperationEnum: %T", src)
	}
	return nil
}

type NullOperationEnum struct {
	OperationEnum OperationEnum
	Valid         bool // Valid is true if OperationEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOperationEnum) Scan(value interface{}) error {
	if value == nil {
		ns.OperationEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OperationEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOperationEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OperationEnum), nil
}

type Account struct {
	ID           int32
	Name         string
	AccountLimit pgtype.Int4
	Balance      pgtype.Int4
}

type Transaction struct {
	ID          int32
	AccountID   pgtype.Int4
	Amount      int32
	Operation   NullOperationEnum
	Description string
	Timestamp   pgtype.Timestamp
}