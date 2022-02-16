package main

import (
	"database/sql"
	"log"

	"github.com/nei7/bank/api"
	"github.com/nei7/bank/internal/db"

	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
	serverAddr = ":3000"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
}
