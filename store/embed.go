package store

import "embed"

var (
	//go:embed migrations/*.sql
	FS embed.FS
)
