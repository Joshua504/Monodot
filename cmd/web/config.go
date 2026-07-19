package main

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	UploadDir       string
	OutputDir       string
	DefaultCellSize int
	MaxUploadSize   int64
}

func NewConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", ":8080"),
		UploadDir:       getEnv("UPLOAD_DIR", "uploads"),
		OutputDir:       getEnv("OUTPUT_DIR", "outputs"),
		DefaultCellSize: getEnvInt("DEFAULT_CELL_SIZE", 3),
		MaxUploadSize:   getEnvInt64("MAX_UPLOAD_SIZE", 10<<20),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return n
}
