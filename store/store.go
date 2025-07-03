package store

import (
	"context"
	"uptimemonitor"
)

type Store interface {
	UserStore
}

type UserStore interface {
	CountUsers(context.Context) (int, error)
	CreateUser(context.Context, uptimemonitor.User) (uptimemonitor.User, error)
}
