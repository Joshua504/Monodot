package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache(templateDir string) (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	root := projectRoot()

	templateDir = filepath.Join(root, templateDir)

	pages, err := filepath.Glob(filepath.Join(templateDir, "pages", "*.html"))
	if err != nil {
		return nil, err
	}

	base := filepath.Join(templateDir, "base.html")

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
