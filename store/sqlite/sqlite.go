package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	db *sql.DB

	*UserStore
	*SessionStore
	*MonitorStore
	*CheckStore
}

func New(dsn string) *Sqlite {
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

	goose.SetBaseFS(FS)
	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	return &Sqlite{
		db:           db,
		UserStore:    NewUserStore(db),
		SessionStore: NewSessionStore(db),
		MonitorStore: NewMonitorStore(db),
		CheckStore:   NewCheckStore(db),
	}
}

func (s *Sqlite) DB() *sql.DB {
	return s.db
}
