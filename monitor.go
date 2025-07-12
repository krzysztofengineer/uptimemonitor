package uptimemonitor

import (
	"fmt"
	"time"
)

type Monitor struct {
	ID        int
	Uuid      string
	Url       string
	CreatedAt time.Time
}

func (m Monitor) URI() string {
	return fmt.Sprintf("/m/%s", m.Uuid)
}
