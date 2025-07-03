package html

import "embed"

var (
	//go:embed *.html
	FS embed.FS
)
