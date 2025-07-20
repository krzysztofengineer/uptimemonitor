package store

import (
	"context"
	"database/sql"
	"time"
	"uptimemonitor"

	"github.com/google/uuid"
)

type CheckStore struct {
	db *sql.DB
}

func NewCheckStore(db *sql.DB) *CheckStore {
	return &CheckStore{db: db}
}

func (s *CheckStore) CreateCheck(ctx context.Context, check uptimemonitor.Check) (uptimemonitor.Check, error) {
	stmt := `INSERT INTO checks(uuid, monitor_id, status_code, response_time_ms, created_at) VALUES(?, ?, ?, ?, ?)`
	uuid := uuid.NewString()
	check.CreatedAt = time.Now()

	res, err := s.db.ExecContext(ctx, stmt, uuid, check.MonitorID, check.StatusCode, check.ResponseTimeMs, check.CreatedAt)
	if err != nil {
		return check, err
	}

	id, _ := res.LastInsertId()

	check.ID = id

	return check, nil
}

func (s *CheckStore) ListChecks(ctx context.Context, monitorID int64, limit int) ([]uptimemonitor.Check, error) {
	stmt := `
		SELECT checks.id, checks.uuid, checks.monitor_id, checks.created_at,
		checks.status_code, checks.response_time_ms,
		monitors.id, monitors.uuid, monitors.url, monitors.created_at
		FROM checks 
		LEFT JOIN monitors ON monitors.id = checks.monitor_id
		WHERE monitor_id = ? 
		ORDER BY checks.id DESC 
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, stmt, monitorID, limit)
	if err != nil {
		return []uptimemonitor.Check{}, err
	}

	defer rows.Close()

	var checks []uptimemonitor.Check

	for rows.Next() {
		var c uptimemonitor.Check

		if err := rows.Scan(
			&c.ID, &c.Uuid, &c.MonitorID, &c.CreatedAt,
			&c.StatusCode, &c.ResponseTimeMs,
			&c.Monitor.ID, &c.Monitor.Uuid, &c.Monitor.Url, &c.Monitor.CreatedAt,
		); err != nil {
			return []uptimemonitor.Check{}, err
		}

		checks = append(checks, c)
	}

	if err = rows.Err(); err != nil {
		return []uptimemonitor.Check{}, err
	}

	return checks, nil
}
