package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/db/utils"
	"testing"

	_ "github.com/lib/pq"
)


var testQueries *Queries // global variable
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	
	// create a test database connection
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// initialize the testQueries variable
	testQueries = New(testDB)

	// run the tests
	os.Exit(m.Run())
}