package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	base := "./templates/base.html"

	for _, page := range pages {

		name := filepath.Base(page)

		tmpl, err := template.ParseFiles(
			base,
			page,
		)

		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}
