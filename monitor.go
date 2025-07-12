package uptimemonitor

import "time"

type Monitor struct {
	ID        int
	Uuid      string
	Url       string
	CreatedAt time.Time
}
