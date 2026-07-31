package main

import (
	"bytes"
	"net/http"
)

func (app *application) render(
	w http.ResponseWriter,
	status int,
	page string,
	data any,
) {
	tmpl, ok := app.templateCache[page]
	if !ok {
		app.logger.Error(
			"Template not found",
			"template", page,
		)

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	var buf bytes.Buffer

	err := tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		app.logger.Error(
			"Template execution failed",
			"template", page,
			"error", err,
		)

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		app.logger.Error(
			"Failed to write response",
			"template", page,
			"error", err,
		)
	}
}
