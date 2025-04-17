package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // Postgres driver
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries // contain DBTX that helps connect to the database

func TestMain(m *testing.M) {
	// create a connection to the database
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}
	// create a testQueries object
	testQueries = New(conn)
	// run the tests and close the connection
	os.Exit(m.Run()) // return an exit call telling us whether the test passed or failed
}
