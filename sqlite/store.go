package sqlite

import (
	"database/sql"
	"uptimemonitor"

	_ "modernc.org/sqlite"
)

type Store struct {
	uptimemonitor.UserStore
}

func New(dsn string) *Store {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Store{
		UserStore: NewUserStore(db),
	}
}
