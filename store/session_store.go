package store

import (
	"context"
	"time"
	"uptimemonitor"

	"github.com/google/uuid"
)

func (s *Store) CreateSession(ctx context.Context, session uptimemonitor.Session) (uptimemonitor.Session, error) {
	stmt := `INSERT INTO sessions (uuid, user_id, created_at, expires_at) VALUES(?, ?, ?, ?)`
	uuid := uuid.NewString()
	session.CreatedAt = time.Now()

	res, err := s.db.ExecContext(ctx, stmt, uuid, session.UserID, session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return session, err
	}

	id, _ := res.LastInsertId()

	session.ID = id
	session.Uuid = uuid
	return session, nil
}

func (s *Store) GetSessionByUuid(ctx context.Context, uuid string) (uptimemonitor.Session, error) {
	stmt := `
		SELECT sessions.id, sessions.user_id, sessions.created_at, sessions.expires_at,
				users.id, users.name, users.email, users.created_at
		FROM sessions 
		LEFT JOIN users ON users.id = sessions.user_id
		WHERE uuid = ? 
		LIMIT 1
	`
	var session uptimemonitor.Session

	err := s.db.QueryRowContext(ctx, stmt, uuid).Scan(
		&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt,
		&session.User.ID, &session.User.Name, &session.User.Email, &session.User.CreatedAt,
	)

	return session, err
}

func (s *Store) RemoveSessionByID(ctx context.Context, id int64) error {
	stmt := `
		DELETE FROM sessions WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, stmt, id)
	return err
}
