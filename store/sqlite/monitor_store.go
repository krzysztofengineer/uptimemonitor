package sqlite

import (
	"context"
	"database/sql"
	"uptimemonitor"

	"github.com/google/uuid"
)

type MonitorStore struct {
	db *sql.DB
}

func NewMonitorStore(db *sql.DB) *MonitorStore {
	return &MonitorStore{db: db}
}

func (s *MonitorStore) CreateMonitor(ctx context.Context, monitor uptimemonitor.Monitor) (uptimemonitor.Monitor, error) {
	stmt := `INSERT INTO monitors(url, uuid, created_at) VALUES(?,?,?)`

	uuid := uuid.NewString()
	_, err := s.db.ExecContext(ctx, stmt, monitor.Url, uuid, monitor.CreatedAt)
	if err != nil {
		return monitor, err
	}

	monitor.Uuid = uuid
	return monitor, nil
}

func (s *MonitorStore) ListMonitors(ctx context.Context) ([]uptimemonitor.Monitor, error) {
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

func (s *MonitorStore) GetMonitorByID(ctx context.Context, id int) (uptimemonitor.Monitor, error) {
	stmt := `SELECT id, url, uuid, created_at FROM monitors WHERE id = ? LIMIT 1`
	var m uptimemonitor.Monitor
	err := s.db.QueryRowContext(ctx, stmt, id).Scan(&m.ID, &m.Url, &m.Uuid, &m.CreatedAt)
	return m, err
}
