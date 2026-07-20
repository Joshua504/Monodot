package main

import "net/http"

func (app *application) serverError(w http.ResponseWriter, err any) {
	app.logger.Printf("Server Error: %v", err)
	http.Error(w, "Something went wrong on our end.", http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) validationError(w http.ResponseWriter, message string) {
	app.logger.Printf("VALIDATION ERROR: %s", message)
	http.Error(w, message, http.StatusBadRequest)
}
