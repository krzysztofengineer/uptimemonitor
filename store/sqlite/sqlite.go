package sqlite

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	db *sql.DB
	*UserStore
	*SessionStore
	*MonitorStore
}

func New(dsn string) *Sqlite {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
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
	}
}

func (s *Sqlite) DB() *sql.DB {
	return s.db
}
