package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/manther/simplebank/db/util"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	// var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: %w", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSourceTest)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
