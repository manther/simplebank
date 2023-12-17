package main

import (
	"database/sql"
	"log"

	"github.com/manther/simplebank/api"
	db "github.com/manther/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

// TODO: Store in secret vault
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	server.Start(":8080")
}
