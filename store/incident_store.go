package store

import (
	"context"
	"time"
	"uptimemonitor"

	"github.com/google/uuid"
)

func (s *Store) CreateIncident(ctx context.Context, incident uptimemonitor.Incident) (uptimemonitor.Incident, error) {
	stmt := `INSERT INTO incidents (uuid, monitor_id, status_text, status_code, response_time_ms, body, headers, req_method, req_url, req_headers, req_body, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	if incident.CreatedAt.IsZero() {
		incident.CreatedAt = time.Now()
	}

	incident.Uuid = uuid.NewString()
	incident.StatusText = uptimemonitor.IncidentStatusOpen

	res, err := s.db.ExecContext(ctx, stmt, incident.Uuid, incident.MonitorID, incident.StatusText, incident.StatusCode, incident.ResponseTimeMs, incident.Body, incident.Headers,
		incident.ReqMethod, incident.ReqUrl, incident.ReqHeaders, incident.ReqBody, incident.CreatedAt)
	if err != nil {
		return uptimemonitor.Incident{}, err
	}

	id, err := res.LastInsertId()
	incident.ID = id

	return incident, err
}

func (s *Store) UpdateIncidentBodyAndHeaders(ctx context.Context, incident uptimemonitor.Incident, body, headers, reqMethod, reqUrl, reqHeaders, reqBody string) error {
	stmt := `
		UPDATE incidents
		SET body = ?, headers = ?, req_method = ?, req_url = ?, req_headers = ?, req_body = ?
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, stmt, body, headers, reqMethod, reqUrl, reqHeaders, reqBody, incident.ID)

	return err
}

func (s *Store) LastIncidentByStatusCode(ctx context.Context, monitorID int64, status string, statusCode int) (uptimemonitor.Incident, error) {
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
		&incident.StatusText, &incident.StatusCode, &incident.ResponseTimeMs,
		&incident.Body, &incident.Headers, &incident.CreatedAt,
	); err != nil {
		return uptimemonitor.Incident{}, err
	}

	return incident, nil
}

func (s *Store) LastOpenIncident(ctx context.Context, monitorID int64) (uptimemonitor.Incident, error) {
	stmt := `
		SELECT id, uuid, monitor_id, status_text, status_code, response_time_ms, body, headers, created_at
		FROM incidents 
		WHERE monitor_id = ? AND status_text = ?
		ORDER BY id DESC 
		LIMIT 1
	`

	row := s.db.QueryRowContext(ctx, stmt, monitorID, uptimemonitor.IncidentStatusOpen)

	var incident uptimemonitor.Incident
	if err := row.Scan(
		&incident.ID, &incident.Uuid, &incident.MonitorID,
		&incident.StatusText, &incident.StatusCode, &incident.ResponseTimeMs,
		&incident.Body, &incident.Headers, &incident.CreatedAt,
	); err != nil {
		return uptimemonitor.Incident{}, err
	}

	return incident, nil
}

func (s *Store) ListOpenIncidents(ctx context.Context) ([]uptimemonitor.Incident, error) {
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
			&incident.StatusText, &incident.StatusCode, &incident.ResponseTimeMs,
			&incident.Body, &incident.Headers, &incident.CreatedAt,
			&incident.Monitor.ID, &incident.Monitor.Url, &incident.Monitor.Uuid, &incident.Monitor.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}

func (s *Store) ListMonitorIncidents(ctx context.Context, id int64) ([]uptimemonitor.Incident, error) {
	stmt := `
		SELECT incidents.id, incidents.uuid, incidents.monitor_id,
			incidents.status_text, incidents.status_code, incidents.response_time_ms,
			incidents.body, incidents.headers, incidents.created_at, incidents.resolved_at,
			monitors.id, monitors.url, monitors.uuid, monitors.created_at
		FROM incidents
		JOIN monitors ON incidents.monitor_id = monitors.id
		WHERE incidents.monitor_id = ?
		ORDER BY incidents.id DESC
		LIMIT 10
	`

	rows, err := s.db.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []uptimemonitor.Incident
	for rows.Next() {
		var incident uptimemonitor.Incident
		if err := rows.Scan(
			&incident.ID, &incident.Uuid, &incident.MonitorID,
			&incident.StatusText, &incident.StatusCode, &incident.ResponseTimeMs,
			&incident.Body, &incident.Headers, &incident.CreatedAt, &incident.ResolvedAt,
			&incident.Monitor.ID, &incident.Monitor.Url, &incident.Monitor.Uuid, &incident.Monitor.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}
func (s *Store) CountMonitorIncidents(ctx context.Context, id int64) int64 {
	stmt := `
		SELECT COUNT(*)
		FROM incidents
		WHERE monitor_id = ?
	`

	var count int64
	s.db.QueryRowContext(ctx, stmt, id).Scan(&count)

	return count
}

func (s *Store) ListMonitorOpenIncidents(ctx context.Context, id int64) ([]uptimemonitor.Incident, error) {
	stmt := `
		SELECT incidents.id, incidents.uuid, incidents.monitor_id,
			incidents.status_text, incidents.status_code, incidents.response_time_ms,
			incidents.body, incidents.headers, incidents.created_at,
			monitors.id, monitors.url, monitors.uuid, monitors.created_at
		FROM incidents
		JOIN monitors ON incidents.monitor_id = monitors.id
		WHERE incidents.monitor_id = ? AND incidents.status_text = ?
		ORDER BY incidents.id DESC
	`

	rows, err := s.db.QueryContext(ctx, stmt, id, uptimemonitor.IncidentStatusOpen)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []uptimemonitor.Incident
	for rows.Next() {
		var incident uptimemonitor.Incident
		if err := rows.Scan(
			&incident.ID, &incident.Uuid, &incident.MonitorID,
			&incident.StatusText, &incident.StatusCode, &incident.ResponseTimeMs,
			&incident.Body, &incident.Headers, &incident.CreatedAt,
			&incident.Monitor.ID, &incident.Monitor.Url, &incident.Monitor.Uuid, &incident.Monitor.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}

func (s *Store) ResolveIncident(ctx context.Context, incident uptimemonitor.Incident) error {
	stmt := `
		UPDATE incidents SET status_text = ?, resolved_at = ? WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, stmt, uptimemonitor.IncidentStatusResolved, time.Now(), incident.ID)

	return err
}

func (s *Store) ResolveMonitorIncidents(ctx context.Context, monitor uptimemonitor.Monitor) error {
	stmt := `
		UPDATE incidents SET status_text = ?, resolved_at = ? WHERE monitor_id = ?
	`

	_, err := s.db.ExecContext(ctx, stmt, uptimemonitor.IncidentStatusResolved, time.Now(), monitor.ID)

	return err
}

func (s *Store) DeleteOldIncidents(ctx context.Context) error {
	stmt := `DELETE FROM incidents WHERE created_at < ?`

	_, err := s.db.ExecContext(ctx, stmt, time.Now().Add(-time.Hour*24*7))

	return err
}

func (s *Store) DeleteIncident(ctx context.Context, id int64) error {
	stmt := `DELETE FROM incidents WHERE id = ?`

	_, err := s.db.ExecContext(ctx, stmt, id)

	return err
}

func (s *Store) GetIncidentByUuid(ctx context.Context, uuid string) (uptimemonitor.Incident, error) {
	stmt := `
		SELECT 
			incidents.id, incidents.uuid, incidents.monitor_id,
			incidents.status_text, incidents.status_code, incidents.response_time_ms,
			incidents.body, incidents.headers, incidents.created_at, incidents.resolved_at,
			incidents.req_method, incidents.req_url, incidents.req_headers, incidents.req_body,
			monitors.id, monitors.url, monitors.uuid, monitors.created_at
		FROM incidents
		LEFT JOIN monitors ON monitors.id = incidents.monitor_id
		WHERE incidents.uuid = ?
	`

	row := s.db.QueryRowContext(ctx, stmt, uuid)

	var incident uptimemonitor.Incident
	if err := row.Scan(
		&incident.ID, &incident.Uuid, &incident.MonitorID,
		&incident.StatusText, &incident.StatusCode, &incident.ResponseTimeMs,
		&incident.Body, &incident.Headers, &incident.CreatedAt, &incident.ResolvedAt,
		&incident.ReqMethod, &incident.ReqUrl, &incident.ReqHeaders, &incident.ReqBody,
		&incident.Monitor.ID, &incident.Monitor.Url, &incident.Monitor.Uuid, &incident.Monitor.CreatedAt,
	); err != nil {
		return incident, err
	}

	return incident, nil
}

func (s *Store) GetIncidentByID(ctx context.Context, id int64) (uptimemonitor.Incident, error) {
	stmt := `
		SELECT 
			incidents.id, incidents.uuid, incidents.monitor_id,
			incidents.status_text, incidents.status_code, incidents.response_time_ms,
			incidents.body, incidents.headers, incidents.created_at, incidents.resolved_at,
			incidents.req_method, incidents.req_url, incidents.req_headers, incidents.req_body,
			monitors.id, monitors.url, monitors.uuid, monitors.created_at
		FROM incidents
		LEFT JOIN monitors ON monitors.id = incidents.monitor_id
		WHERE incidents.id = ?
	`

	row := s.db.QueryRowContext(ctx, stmt, id)

	var incident uptimemonitor.Incident
	if err := row.Scan(
		&incident.ID, &incident.Uuid, &incident.MonitorID,
		&incident.StatusText, &incident.StatusCode, &incident.ResponseTimeMs,
		&incident.Body, &incident.Headers, &incident.CreatedAt, &incident.ResolvedAt,
		&incident.ReqMethod, &incident.ReqUrl, &incident.ReqHeaders, &incident.ReqBody,
		&incident.Monitor.ID, &incident.Monitor.Url, &incident.Monitor.Uuid, &incident.Monitor.CreatedAt,
	); err != nil {
		return incident, err
	}

	return incident, nil
}
