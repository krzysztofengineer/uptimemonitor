package app

import (
	"uptimemonitor"

	_ "modernc.org/sqlite"
)

type Store interface {
	uptimemonitor.UserStore
}
