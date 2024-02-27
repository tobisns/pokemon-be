package db

import (
	"database/sql"
	"fmt" // postgres driver side effects for migrations

	_ "github.com/lib/pq"
)

// GetConnection ...
func GetConnection(host string, port int, user, password, dbName string) (*sql.DB, error) {
	if password == "" { // local DBs my not require a password
		password = `''`
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
