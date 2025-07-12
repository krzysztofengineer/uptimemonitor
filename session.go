package uptimemonitor

import "time"

type Session struct {
	ID        int
	UserID    int
	Uuid      string
	CreatedAt time.Time
	ExpiresAt time.Time
	User      User
}
