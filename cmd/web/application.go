package main

import (
	"html/template"
	"log/slog"
)

type application struct {
	config         *Config
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	imageGenerator func(string, string, int) error
}
