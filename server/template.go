package server

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
		data, _, err := fetch(r)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		err = h.tmpl[templateFile].ExecuteTemplate(w, templateEntry, data)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		w.Header().Set("content-type", "text/html")
	}
}
