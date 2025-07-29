package uptimemonitor

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Monitor struct {
	ID                int64
	Uuid              string
	Url               string
	HttpMethod        string
	HttpHeaders       string
	HttpBody          string
	WebhookUrl        string
	WebhookMethod     string
	WebhookHeaders    string
	WebhookBody       string
	Uptime            float32
	AvgResponseTimeMs int64
	N                 int64
	IncidentsCount    int64
	CreatedAt         time.Time
}

func (m Monitor) URI() string {
	return fmt.Sprintf("/m/%s", m.Uuid)
}

func (m Monitor) Secure() bool {
	uri, err := url.ParseRequestURI(m.Url)
	if err != nil {
		return false
	}

	return uri.Scheme == "https"
}

func (m Monitor) Domain() string {
	uri, err := url.ParseRequestURI(m.Url)
	if err != nil {
		return m.Url
	}

	res, err := url.JoinPath(uri.Host, strings.TrimLeft(uri.Path, "/"))
	if err != nil {
		return m.Url
	}

	if uri.RawQuery != "" {
		res = fmt.Sprintf("%s?%s", res, uri.Query().Encode())
	}

	return res
}

func (m Monitor) Host() string {
	uri, err := url.ParseRequestURI(m.Url)
	if err != nil {
		return m.Url
	}

	host := uri.Host
	maxLimit := 30

	if len(host) > maxLimit {
		return fmt.Sprintf("%s...", host[:maxLimit])
	}

	return host
}
