package sqlite

import "embed"

var (
	//go:embed migrations/*.sql
	FS embed.FS
)
