package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/Bronku/cake-order-tracker/logging"
)

//go:embed templates/*
var templates embed.FS

func (h *Server) loadTemplates() {
	h.tmpl = make(map[string]*template.Template)
	for _, page := range h.routes {
		tmpl, err := template.ParseFS(templates,
			"templates/layout/*.gohtml",
			"templates/"+page.template+".gohtml",
		)
		if err != nil {
			log.Fatal(err)
		}
		h.tmpl[page.template] = tmpl
	}
}

func (h *Server) render(fetch fetcher, templateFile string, templateEntry string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, code, err := fetch(r)
		if err != nil {
			logging.ErrorPage(err, code).ServeHTTP(w, r)
			return
		}
		err = h.tmpl[templateFile].ExecuteTemplate(w, templateEntry, data)
		if err != nil {
			logging.ErrorPage(err, http.StatusInternalServerError).ServeHTTP(w, r)
			return
		}
		w.Header().Set("content-type", "text/html")
	}
}
