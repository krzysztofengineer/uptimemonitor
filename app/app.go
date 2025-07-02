package app

import "database/sql"

type App struct {
	Store *Store
}

func New(db *sql.DB) *App {
	return &App{
		Store: NewStore(db),
	}
}
