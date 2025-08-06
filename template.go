package main

import (
	"embed"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/layout/*
var layout embed.FS

//go:embed templates/pages/*
var pages embed.FS

var templates map[string]*template.Template

func loadTemplates() {
	// Read all page template files
	pageFiles, err := pages.ReadDir("templates/pages")
	if err != nil {
		log.Fatal(err)
	}

	for _, pageFile := range pageFiles {
		if pageFile.IsDir() {
			continue
		}

		// Extract template name from filename (without extension)
		templateName := strings.TrimSuffix(pageFile.Name(), filepath.Ext(pageFile.Name()))

		// Create template with layout and specific page
		tmpl, err := template.ParseFS(layout, "templates/layout/*.gohtml")
		if err != nil {
			log.Fatal(err)
		}

		// Parse the specific page template
		pagePattern := "templates/pages/" + pageFile.Name()
		tmpl, err = tmpl.ParseFS(pages, pagePattern)
		if err != nil {
			log.Fatal(err)
		}

		templates[templateName] = tmpl
	}
}
