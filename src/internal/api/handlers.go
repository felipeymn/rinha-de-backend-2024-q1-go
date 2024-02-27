package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) SaveTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := getUrlId(ps)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	transaction, err := parsePostTransactionRequest(w, r)
	if err != nil {
		return
	}

	accountRow, err := s.storage.SaveTransaction(id, transaction)
	if err != nil {
		if err.Error() == ErrNoRows.Error() {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	// Create account response
	account := transactionDAOToDTO(accountRow)
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(&account)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(response))
}

func (s *Server) GetStatement(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := getUrlId(ps)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	accountRow, transactionRows, err := s.storage.GetStatement(id)
	if err != nil {
		if err.Error() == ErrNoRows.Error() {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	statement := statementDAOToDTO(accountRow, transactionRows)

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(&statement)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(response))
}
