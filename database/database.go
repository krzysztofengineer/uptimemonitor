package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Must(db *sql.DB, err error) *sql.DB {
	if err != nil {
		panic(err)
	}
	return db
}
