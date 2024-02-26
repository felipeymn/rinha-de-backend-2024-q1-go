package main

import (
	"log"

	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/api"
	"github.com/felipeymn/rinha-de-backend-2024-q1/src/internal/database"
)

func main() {
	// start db connection pool
	db := database.NewPostgreSQL()
	defer db.Pool.Close()

	// start http server
	server := api.NewServer(":8080", db)
	log.Fatal(server.Start())
}
