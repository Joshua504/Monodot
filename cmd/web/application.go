package main

import (
	"html/template"
	"log"
)

type application struct {
	config        *Config
	logger        *log.Logger
	templateCache map[string]*template.Template
}
