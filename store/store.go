package store

import (
	"context"
	"database/sql"
	"uptimemonitor"
)

type Store interface {
	DB() *sql.DB
	UserStore
}

type UserStore interface {
	CountUsers(context.Context) (int, error)
	CreateUser(context.Context, uptimemonitor.User) (uptimemonitor.User, error)
}
