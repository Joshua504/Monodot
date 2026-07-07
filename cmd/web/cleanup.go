package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func cleanupDirectory(dir string, maxAge time.Duration) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(dir, file.Name())

		info, err := file.Info()
		if err != nil {
			continue
		}

		if time.Since(info.ModTime()) > maxAge {
			err := os.Remove(path)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("deleted: %s\n", path)
		}
	}

}

func startCleanupJob() {
	cleanupDirectory("uploads", time.Hour)
	cleanupDirectory("outputs", time.Hour)

	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for {
			<-ticker.C

			cleanupDirectory(
				"uploads",
				time.Hour,
			)

			cleanupDirectory(
				"outputs",
				time.Hour,
			)
		}
	}()
}
