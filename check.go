package uptimemonitor

import "time"

type Check struct {
	ID             int64
	Uuid           string
	MonitorID      int64
	StatusCode     int
	ResponseTimeMs int64
	CreatedAt      time.Time

	Monitor Monitor
}

func (c Check) ColorClass() string {
	if c.StatusCode >= 200 && c.StatusCode < 300 {
		return "bg-lime-300"
	} else if c.StatusCode >= 300 && c.StatusCode < 400 {
		return "bg-yellow-300"
	} else if c.StatusCode >= 400 && c.StatusCode < 500 {
		return "bg-orange-300"
	} else if c.StatusCode >= 500 {
		return "bg-red-400"
	} else {
		return "bg-neutral-300"
	}
}
