package uptimemonitor

import "time"

const (
	IncidentStatusOpen     = "open"
	IncidentStatusResolved = "resolved"
)

type Incident struct {
	ID             int64
	Uuid           string
	MonitorID      int64
	CreatedAt      time.Time
	StatusCode     int
	ResponseTimeMs int64
	Body           string
	Headers        string
	Status         string

	Monitor Monitor
}
