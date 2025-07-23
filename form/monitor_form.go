package form

import (
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
)

type MonitorForm struct {
	Url              string
	HttpMethod       string
	HasCustomHeaders bool
	HttpHeaders      string
	HasCustomBody    bool
	HttpBody         string

	Errors map[string]string
}

func (f *MonitorForm) Validate() bool {
	f.Errors = map[string]string{}

	if f.Url == "" {
		f.Errors["Url"] = "The url is required"
	} else if _, err := url.ParseRequestURI(f.Url); err != nil {
		f.Errors["Url"] = "The url is invalid"
	}

	methods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete,
	}

	if !slices.Contains(methods, f.HttpMethod) {
		f.Errors["HttpMethod"] = "The http method is invalid"
	}

	if f.HasCustomHeaders {
		headers := map[string]any{}
		err := json.Unmarshal([]byte(f.HttpHeaders), &headers)

		if err != nil {
			f.Errors["HttpHeaders"] = "The http headers should be a valid JSON"
		}
	}

	return len(f.Errors) == 0
}
