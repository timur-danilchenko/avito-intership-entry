package database

import (
	"database/sql"
	"fmt"
	"timur-danilchenko/avito-intership-entry/src/main/utilities"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	dbUser := utilities.GetEnv("DB_USER", "postgres")
	dbPassword := utilities.GetEnv("DB_PASS", "postgres")
	dbHost := utilities.GetEnv("DB_HOST", "localhost")
	dbPort := utilities.GetEnv("DB_PORT", "5432")
	dbName := utilities.GetEnv("DB_NAME", "postgres")

	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
