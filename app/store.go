package app

import (
	"database/sql"
	"uptimemonitor/store"
)

type Store struct {
	*store.UserStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		UserStore: store.NewUserStore(db),
	}
}
