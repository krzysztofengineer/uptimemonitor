package store

import (
	"context"
	"database/sql"
	"uptimemonitor"
)

type Store interface {
	DB() *sql.DB
	UserStore
	SessionStore
	MonitorStore
}

type UserStore interface {
	CountUsers(context.Context) (int, error)
	CreateUser(context.Context, uptimemonitor.User) (uptimemonitor.User, error)
	GetUserByEmail(context.Context, string) (uptimemonitor.User, error)
}

type SessionStore interface {
	CreateSession(context.Context, uptimemonitor.Session) (uptimemonitor.Session, error)
	GetSessionByUuid(context.Context, string) (uptimemonitor.Session, error)
}

type MonitorStore interface {
	CreateMonitor(context.Context, uptimemonitor.Monitor) (uptimemonitor.Monitor, error)
	ListMonitors(context.Context) ([]uptimemonitor.Monitor, error)
}
