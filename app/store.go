package app

import (
	"database/sql"
	"uptimemonitor"
	"uptimemonitor/store"

	_ "modernc.org/sqlite"
)

type Store struct {
	uptimemonitor.UserStore
}

func MustNewStore(dsn string) *Store {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Store{
		UserStore: store.NewUserStore(db),
	}
}
