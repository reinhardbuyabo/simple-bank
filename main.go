package main

import (
	"database/sql"
	"log"

	"github.com/reinhardbuyabo/simplebank/api"
	db "github.com/reinhardbuyabo/simplebank/db/sqlc"

	_ "github.com/lib/pq" // postgres driver
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8081"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
