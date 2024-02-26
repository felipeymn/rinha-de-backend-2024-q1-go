package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/database/sqlc"
	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/types"
	"github.com/julienschmidt/httprouter"
)

var ErrNoRows = errors.New("no rows in result set")

func getUrlId(ps httprouter.Params) (id int32, err error) {
	idParam := ps.ByName("id")
	// Convert string to int32
	num, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return
	}
	// Convert int64 to int32
	id = int32(num)
	return
}

func statementDAOToDTO(
	accountRow sqlc.GetAccountRow,
	transactionRows []sqlc.GetTransactionsRow,
) (statement types.StatementResponse) {
	// set account information
	statement.Account = types.StatementAccount{
		Limit:     accountRow.AccountLimit,
		Balance:   accountRow.Balance,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// set transactions
	var statementTransactions []types.StatementTransaction
	for _, transactionRow := range transactionRows {
		statementTransactions = append(statementTransactions, types.StatementTransaction{
			Amount:      transactionRow.Amount,
			Operation:   transactionRow.Operation.OperationEnum,
			Description: transactionRow.Description,
			Timestamp:   transactionRow.Timestamp})
	}
	statement.Transactions = statementTransactions

	return statement
}

func transactionDAOToDTO(
	account sqlc.UpdateAccountRow,
) (transactionResponse types.TransactionResponse) {
	transactionResponse.Limit = account.AccountLimit
	transactionResponse.Balance = account.Balance.Int32
	return transactionResponse
}

func parsePostTransactionRequest(
	w http.ResponseWriter,
	r *http.Request,
) (transaction types.Transaction, err error) {
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return transaction, err
	}
	defer r.Body.Close()

	if err := transaction.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return transaction, err
	}
	return transaction, nil
}
