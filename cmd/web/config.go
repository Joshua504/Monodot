package main

type Config struct {
	Port            string
	UploadDir       string
	OutputDir       string
	DefaultCellSize int
	MaxUploadSize   int64
}

func NewConfig() *Config {
	return &Config{
		Port:            ":8080",
		UploadDir:       "uploads",
		OutputDir:       "outputs",
		DefaultCellSize: 3,
		MaxUploadSize:   10 << 20,
	}
}
