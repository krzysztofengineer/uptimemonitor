package uptimemonitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	IncidentStatusOpen     string = "open"
	IncidentStatusResolved string = "resolved"
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
	StatusText     string
	ResolvedAt     *time.Time
	ReqMethod      string
	ReqUrl         string
	ReqHeaders     string
	ReqBody        string

	Monitor Monitor
}

func (i Incident) URI() string {
	return fmt.Sprintf("/m/%s/i/%s", i.Monitor.Uuid, i.Uuid)
}

func (i Incident) BadgeClass() string {
	if i.StatusCode >= 200 && i.StatusCode < 300 {
		return "badge-success"
	} else if i.StatusCode >= 300 && i.StatusCode < 400 {
		return "badge-warning"
	} else if i.StatusCode >= 400 && i.StatusCode < 500 {
		return "badge-accent"
	} else if i.StatusCode >= 500 {
		return "badge-error"
	} else {
		return "badge-neutral"
	}
}

func (i Incident) StatusCodeText() string {
	return http.StatusText(i.StatusCode)
}

func (i Incident) StatusBadgeClass() string {
	if i.StatusText == IncidentStatusOpen {
		return "badge-error"
	}

	return "badge-success"
}

func (i Incident) StatusBadgeText() string {
	if i.StatusText == IncidentStatusResolved {
		return "Resolved"
	}

	return "Open"
}

func (i Incident) ReqHeadersMap() map[string]string {
	if i.ReqHeaders == "" {
		return map[string]string{}
	}

	customHeaders := map[string]string{}
	err := json.Unmarshal([]byte(i.ReqHeaders), &customHeaders)
	if err != nil {
		return map[string]string{}

	}

	return customHeaders
}
