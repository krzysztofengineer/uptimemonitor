package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Store struct {
	*UserStore
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
