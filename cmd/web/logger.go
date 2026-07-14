package main

import (
	"log"
	"net/http"
	"os"
)

func requestLogger(r *http.Request) *log.Logger {
	requestID := getRequestID(r)
	prefix := "[" + requestID + "]"
	return log.New(os.Stdout, prefix, log.LstdFlags)
}
