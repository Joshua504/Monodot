package main

import (
	"net/http"
)

func (app *application) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {

				app.logger.Error(
					"Recovered panic",
					"error", err,
				)

				app.serverError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
