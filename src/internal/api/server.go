package api

import (
	"net/http"

	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/database"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	listenAddr string
	storage    database.Storage
}

func NewServer(listenAddr string, storage database.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *Server) Start() error {
	router := httprouter.New()
	router.POST("/clientes/:id/transacoes", s.SaveTransaction)
	router.GET("/clientes/:id/extrato", s.GetStatement)

	return http.ListenAndServe(s.listenAddr, router)
}
