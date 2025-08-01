package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor"
	"uptimemonitor/html"
)

func (*Handler) ListSponsors() http.HandlerFunc {
	layout := template.Must(template.ParseFS(html.FS, "layout.html"))
	sponsor := template.Must(template.ParseFS(html.FS, "sponsor.html"))

	type data struct {
		Sponsors []uptimemonitor.Sponsor
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			layout.ExecuteTemplate(w, "sponsors", nil)
			return
		}

		d := data{
			Sponsors: []uptimemonitor.Sponsor{
				{
					Name:  "AIR Labs",
					Url:   "https://airlabs.pl",
					Image: "/static/img/airlabs.svg",
				},
			},
		}

		sponsor.Execute(w, d)
	}
}
