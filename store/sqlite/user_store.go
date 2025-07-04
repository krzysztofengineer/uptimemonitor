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
	stmt := `SELECT COUNT(*) FROM users`

	var count int
	if err := s.db.QueryRowContext(ctx, stmt).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *UserStore) CreateUser(ctx context.Context, user uptimemonitor.User) (uptimemonitor.User, error) {
	stmt := `INSERT INTO users (name, email, password_hash, created_at) VALUES (?, ?, ?, ?) RETURNING id`

	row := s.db.QueryRowContext(ctx, stmt, user.Name, user.Email, user.PasswordHash, user.CreatedAt)

	var id int
	if err := row.Scan(&id); err != nil {
		return uptimemonitor.User{}, err
	}

	user.ID = id

	return user, nil
}
