package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var Dbt *sql.DB
var Qs *Queries

func Setup() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Could not load '.env' config file -- %s", err)
	}

	Dbt, err := sql.Open(
		os.Getenv("DB_ENGINE"),
		os.Getenv("DATABASE"))

	if err != nil {
		log.Fatalf("Could not build database tool -- %s", err)
	}

	err = Dbt.Ping()

	if err != nil {
		log.Fatalf("Could not connect to database -- %s", err)
	}

	Qs = New(Dbt)
}
