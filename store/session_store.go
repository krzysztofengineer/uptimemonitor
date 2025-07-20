package store

import (
	"context"
	"database/sql"
	"time"
	"uptimemonitor"

	"github.com/google/uuid"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) CreateSession(ctx context.Context, session uptimemonitor.Session) (uptimemonitor.Session, error) {
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

func (s *SessionStore) GetSessionByUuid(ctx context.Context, uuid string) (uptimemonitor.Session, error) {
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
