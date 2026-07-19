package main

import (
	"log"
	"net/http"
)

func newServer(cfg *Config) *http.Server {
	app := &application{config: cfg}
	mux := http.NewServeMux()

	mux.Handle(
		"/outputs/",
		http.StripPrefix(
			"/outputs/",
			http.FileServer(http.Dir(cfg.OutputDir)),
		),
	)

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/generate", app.generateHandler)
	mux.HandleFunc("/result", app.resultHandler)

	return &http.Server{
		Addr:    cfg.Port,
		Handler: loggingMiddleware(mux),
	}
}

func startServer(server *http.Server) {
	log.Print("STARTING AND LISTENING TO SERVER ON____", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
