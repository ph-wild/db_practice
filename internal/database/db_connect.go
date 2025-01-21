package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDB(connString string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	return db
}
