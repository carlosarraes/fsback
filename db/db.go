package db

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		return nil, err
	}

	return db, nil
}
