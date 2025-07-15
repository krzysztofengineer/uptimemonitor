package uptimemonitor

import "time"

type Check struct {
	ID        int64
	Uuid      string
	MonitorID int64
	CreatedAt time.Time

	Monitor Monitor
}
