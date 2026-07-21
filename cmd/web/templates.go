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
		app.logger.Printf("Template %q does not exist", page)

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
		app.logger.Printf(
			"Template execute error: %v",
			err,
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
		app.logger.Printf(
			"Response write error: %v",
			err,
		)
	}
}
