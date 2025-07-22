package store

import (
	"context"
	"time"
	"uptimemonitor"
)

func (s *Store) CountUsers(ctx context.Context) (int, error) {
	stmt := `SELECT COUNT(*) FROM users`

	var count int
	if err := s.db.QueryRowContext(ctx, stmt).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CreateUser(ctx context.Context, user uptimemonitor.User) (uptimemonitor.User, error) {
	stmt := `INSERT INTO users (name, email, password_hash, created_at) VALUES (?, ?, ?, ?) RETURNING id`
	user.CreatedAt = time.Now()

	res, err := s.db.ExecContext(ctx, stmt, user.Name, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return uptimemonitor.User{}, err
	}

	id, _ := res.LastInsertId()

	user.ID = id

	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (uptimemonitor.User, error) {
	stmt := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = ? LIMIT 1`

	row := s.db.QueryRowContext(ctx, stmt, email)

	var user uptimemonitor.User
	if err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt,
	); err != nil {
		return user, err
	}

	return user, nil
}
