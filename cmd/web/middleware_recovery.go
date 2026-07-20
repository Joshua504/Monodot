package main

import (
	"net/http"
	"runtime/debug"
)

func (app *application) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {

				app.logger.Printf(
					"PANIC: %v\n%s",
					err,
					debug.Stack(),
				)

				app.serverError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
