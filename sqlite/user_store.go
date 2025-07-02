package sqlite

import (
	"context"
	"database/sql"
	"uptimemonitor"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) CountUsers(ctx context.Context) (int, error) {
	return 0, nil
}

func (s *UserStore) CreateUser(ctx context.Context, user uptimemonitor.User) (uptimemonitor.User, error) {
	return uptimemonitor.User{}, nil
}
