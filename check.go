package uptimemonitor

import (
	"net/http"
	"time"
)

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

func (c Check) HeightClass(maxTime int64) string {
	height := int64(c.ResponseTimeMs) * 100 / maxTime
	if height < 10 {
		return "h-[10%]"
	} else if height < 20 {
		return "h-[20%]"
	} else if height < 30 {
		return "h-[30%]"
	} else if height < 40 {
		return "h-[40%]"
	} else if height < 50 {
		return "h-[50%]"
	} else if height < 60 {
		return "h-[60%]"
	} else if height < 70 {
		return "h-[70%]"
	} else if height < 80 {
		return "h-[80%]"
	} else if height < 90 {
		return "h-[90%]"
	} else if height < 100 {
		return "h-full"
	} else {
		return "h-full"
	}
}

func (c Check) BadgeClass() string {
	if c.StatusCode >= 200 && c.StatusCode < 300 {
		return "badge-success"
	} else if c.StatusCode >= 300 && c.StatusCode < 400 {
		return "badge-warning"
	} else if c.StatusCode >= 400 && c.StatusCode < 500 {
		return "badge-accent"
	} else if c.StatusCode >= 500 {
		return "badge-error"
	} else {
		return "badge-neutral"
	}
}

func (c Check) StatusText() string {
	return http.StatusText(c.StatusCode)
}
