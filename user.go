package uptimemonitor

import (
	"context"
	"time"
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

type UserStore interface {
	CountUsers(context.Context) (int, error)
	CreateUser(context.Context, User) (User, error)
}
