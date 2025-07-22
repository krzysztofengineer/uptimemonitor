package store

import (
	"context"
	"time"
	"uptimemonitor"

	"github.com/google/uuid"
)

func (s *Store) CountMonitors(ctx context.Context) int {
	stmt := `SELECT COUNT(*) FROM monitors`

	var count int
	s.db.QueryRowContext(ctx, stmt).Scan(&count)

	return count
}

func (s *Store) CreateMonitor(ctx context.Context, monitor uptimemonitor.Monitor) (uptimemonitor.Monitor, error) {
	stmt := `INSERT INTO monitors(url, uuid, http_method, created_at) VALUES(?,?,?,?)`
	monitor.CreatedAt = time.Now()

	uuid := uuid.NewString()
	res, err := s.db.ExecContext(ctx, stmt, monitor.Url, uuid, monitor.HttpMethod, monitor.CreatedAt)
	if err != nil {
		return monitor, err
	}

	id, _ := res.LastInsertId()

	monitor.ID = id
	monitor.Uuid = uuid
	return monitor, nil
}

func (s *Store) ListMonitors(ctx context.Context) ([]uptimemonitor.Monitor, error) {
	stmt := "SELECT id, url, uuid, created_at FROM monitors ORDER BY created_at DESC"

	rows, err := s.db.QueryContext(ctx, stmt)
	if err != nil {
		return []uptimemonitor.Monitor{}, err
	}
	defer rows.Close()

	var monitors []uptimemonitor.Monitor

	for rows.Next() {
		var m uptimemonitor.Monitor
		if err := rows.Scan(&m.ID, &m.Url, &m.Uuid, &m.CreatedAt); err != nil {
			return monitors, err
		}

		monitors = append(monitors, m)
	}

	if err = rows.Err(); err != nil {
		return monitors, err
	}

	return monitors, nil
}

func (s *Store) GetMonitorByID(ctx context.Context, id int) (uptimemonitor.Monitor, error) {
	stmt := `SELECT id, url, uuid, http_method, created_at FROM monitors WHERE id = ? LIMIT 1`
	var m uptimemonitor.Monitor
	err := s.db.QueryRowContext(ctx, stmt, id).Scan(&m.ID, &m.Url, &m.Uuid, &m.HttpMethod, &m.CreatedAt)
	return m, err
}

func (s *Store) GetMonitorByUuid(ctx context.Context, uuid string) (uptimemonitor.Monitor, error) {
	stmt := `SELECT id, url, uuid, http_method, created_at FROM monitors WHERE uuid = ? LIMIT 1`
	var m uptimemonitor.Monitor
	err := s.db.QueryRowContext(ctx, stmt, uuid).Scan(&m.ID, &m.Url, &m.Uuid, &m.HttpMethod, &m.CreatedAt)
	return m, err
}
