package uptimemonitor

import (
	"fmt"
	"net/url"
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

func (m Monitor) Domain() string {
	uri, err := url.ParseRequestURI(m.Url)
	if err != nil {
		return m.Url
	}

	res, err := url.JoinPath(uri.Host, uri.Path)
	if err != nil {
		return uri.Host
	}

	return res
}
