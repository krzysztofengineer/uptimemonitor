package store

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB

	*UserStore
	*SessionStore
	*MonitorStore
	*CheckStore
	*IncidentStore
}

func New(dsn string) *Store {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		panic(fmt.Sprintf("Failed to enable WAL mode: %v", err))
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	goose.SetBaseFS(FS)
	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	return &Store{
		db:            db,
		UserStore:     NewUserStore(db),
		SessionStore:  NewSessionStore(db),
		MonitorStore:  NewMonitorStore(db),
		CheckStore:    NewCheckStore(db),
		IncidentStore: NewIncidentStore(db),
	}
}

func (s *Store) DB() *sql.DB {
	return s.db
}
