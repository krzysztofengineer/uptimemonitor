package store

import (
	"context"
	"fmt"
	"time"
	"uptimemonitor"

	"github.com/google/uuid"
)

func (s *Store) CreateCheck(ctx context.Context, check uptimemonitor.Check) (uptimemonitor.Check, error) {
	stmt := `INSERT INTO checks(uuid, monitor_id, status_code, response_time_ms, created_at) VALUES(?, ?, ?, ?, ?)`
	uuid := uuid.NewString()
	if check.CreatedAt.IsZero() {
		check.CreatedAt = time.Now()
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return check, err
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, stmt, uuid, check.MonitorID, check.StatusCode, check.ResponseTimeMs, check.CreatedAt)
	if err != nil {
		return check, err
	}

	id, _ := res.LastInsertId()
	check.ID = id

	stmt = `
		SELECT uptime, avg_response_time_ms, n, incidents_count
		FROM monitors 
		WHERE id = ?
	`
	var n int64
	var uptime float32
	var avgResponseTimeMs int64
	var incidentsCount int64
	err = tx.QueryRowContext(ctx, stmt, check.MonitorID).Scan(&uptime, &avgResponseTimeMs, &n, &incidentsCount)
	if err != nil {
		return check, err
	}

	if check.StatusCode >= 300 {
		incidentsCount++
	}

	stmt = `
		UPDATE monitors 
		SET uptime = ?, avg_response_time_ms = ?, n = ?, incidents_count = ?
		WHERE id = ?
	`
	newIncidentCount := incidentsCount
	newN := n + 1
	newUptime := fmt.Sprintf("%.1f", float32(float32(newN-newIncidentCount)/float32(newN)*float32(100)))
	newAvgResponseTimeMs := (avgResponseTimeMs*n + check.ResponseTimeMs) / newN

	_, err = tx.ExecContext(ctx, stmt, newUptime, newAvgResponseTimeMs, newN, newIncidentCount, check.MonitorID)
	if err != nil {
		return check, err
	}

	tx.Commit()

	return check, nil
}

func (s *Store) ListChecks(ctx context.Context, monitorID int64, limit int) ([]uptimemonitor.Check, error) {
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

func (s *Store) GetCheckByID(ctx context.Context, id int64) (uptimemonitor.Check, error) {
	stmt := `
		SELECT 
		checks.id, checks.uuid, checks.monitor_id, checks.status_code, checks.response_time_ms, checks.created_at,
		monitors.id, monitors.url, monitors.uuid, monitors.http_method, monitors.http_headers, monitors.http_body, monitors.created_at
		FROM checks
		LEFT JOIN monitors ON monitors.id = checks.monitor_id
		WHERE checks.id = ?
	`

	var ch uptimemonitor.Check
	err := s.db.QueryRowContext(ctx, stmt, id).Scan(
		&ch.ID, &ch.Uuid, &ch.MonitorID, &ch.StatusCode, &ch.ResponseTimeMs, &ch.CreatedAt,
		&ch.Monitor.ID, &ch.Monitor.Url, &ch.Monitor.Uuid, &ch.Monitor.HttpMethod, &ch.Monitor.HttpHeaders, &ch.Monitor.HttpBody, &ch.Monitor.CreatedAt,
	)

	return ch, err
}

func (s *Store) DeleteOldChecks(ctx context.Context) error {
	stmt := `DELETE FROM checks WHERE created_at < ?`

	_, err := s.db.ExecContext(ctx, stmt, time.Now().Add(-time.Hour))

	return err

}
