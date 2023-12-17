package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/manther/simplebank/api"
	db "github.com/manther/simplebank/db/sqlc"
	"github.com/manther/simplebank/db/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server at: %s, err: %v", config.ServerAddress, err)
	}
}
