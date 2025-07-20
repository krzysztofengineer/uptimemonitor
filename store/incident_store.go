package store

import (
	"context"
	"database/sql"
	"time"
	"uptimemonitor"

	"github.com/google/uuid"
)

type IncidentStore struct {
	db *sql.DB
}

func NewIncidentStore(db *sql.DB) *IncidentStore {
	return &IncidentStore{db: db}
}

func (s *IncidentStore) CreateIncident(ctx context.Context, incident uptimemonitor.Incident) (uptimemonitor.Incident, error) {
	stmt := `INSERT INTO incidents (uuid, monitor_id, status_text, status_code, response_time_ms, body, headers, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	incident.CreatedAt = time.Now()
	incident.Uuid = uuid.NewString()
	incident.Status = uptimemonitor.IncidentStatusOpen

	res, err := s.db.ExecContext(ctx, stmt, incident.Uuid, incident.MonitorID, incident.Status, incident.StatusCode, incident.ResponseTimeMs, incident.Body, incident.Headers, incident.CreatedAt)
	if err != nil {
		return uptimemonitor.Incident{}, err
	}

	id, err := res.LastInsertId()
	incident.ID = id

	return incident, err
}

func (s *IncidentStore) LastIncident(ctx context.Context, monitorID int64, status string, statusCode int) (uptimemonitor.Incident, error) {
	stmt := `
		SELECT id, uuid, monitor_id, status_text, status_code, response_time_ms, body, headers, created_at
		FROM incidents 
		WHERE monitor_id = ? AND status_text = ? AND status_code = ?
		ORDER BY id DESC 
		LIMIT 1
	`

	row := s.db.QueryRowContext(ctx, stmt, monitorID, status, statusCode)

	var incident uptimemonitor.Incident
	if err := row.Scan(
		&incident.ID, &incident.Uuid, &incident.MonitorID,
		&incident.Status, &incident.StatusCode, &incident.ResponseTimeMs,
		&incident.Body, &incident.Headers, &incident.CreatedAt,
	); err != nil {
		return uptimemonitor.Incident{}, err
	}

	return incident, nil
}

func (s *IncidentStore) ListOpenIncidents(ctx context.Context) ([]uptimemonitor.Incident, error) {
	stmt := `
		SELECT incidents.id, incidents.uuid, incidents.monitor_id,
			incidents.status_text, incidents.status_code, incidents.response_time_ms,
			incidents.body, incidents.headers, incidents.created_at,
			monitors.id, monitors.url, monitors.uuid, monitors.created_at
		FROM incidents
		JOIN monitors ON incidents.monitor_id = monitors.id
		WHERE incidents.status_text = ?
		ORDER BY incidents.id DESC
	`

	rows, err := s.db.QueryContext(ctx, stmt, uptimemonitor.IncidentStatusOpen)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []uptimemonitor.Incident
	for rows.Next() {
		var incident uptimemonitor.Incident
		if err := rows.Scan(
			&incident.ID, &incident.Uuid, &incident.MonitorID,
			&incident.Status, &incident.StatusCode, &incident.ResponseTimeMs,
			&incident.Body, &incident.Headers, &incident.CreatedAt,
			&incident.Monitor.ID, &incident.Monitor.Url, &incident.Monitor.Uuid, &incident.Monitor.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}
