package main

import (
	"embed"
	"log"
	"text/template"
)

//go:embed templates/layout/*
var layout embed.FS

//go:embed templates/pages/*
var pages embed.FS

var templates map[string]*template.Template

func loadTemplates() {
	templates = make(map[string]*template.Template)
	pageFiles, err := pages.ReadDir("templates/pages")
	if err != nil {
		log.Fatal(err)
	}

	for _, pageFile := range pageFiles {
		if pageFile.IsDir() {
			continue
		}

		tmpl, err := template.ParseFS(layout, "templates/layout/*.gohtml")
		if err != nil {
			log.Fatal(err)
		}

		pagePattern := "templates/pages/" + pageFile.Name()
		tmpl, err = tmpl.ParseFS(pages, pagePattern)
		if err != nil {
			log.Fatal(err)
		}

		templates[pageFile.Name()] = tmpl
	}
}
