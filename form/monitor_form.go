package form

import "net/url"

type MonitorForm struct {
	Url string

	Errors map[string]string
}

func (f *MonitorForm) Validate() bool {
	f.Errors = map[string]string{}

	if f.Url == "" {
		f.Errors["Url"] = "The url is required"
	} else if _, err := url.ParseRequestURI(f.Url); err != nil {
		f.Errors["Url"] = "The url is invalid"
	}

	return len(f.Errors) == 0
}
