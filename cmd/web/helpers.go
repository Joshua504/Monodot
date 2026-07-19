package main

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func validateExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png":
		return nil
	default:
		return errors.New("File not supported")
	}
}

func validateContentType(file multipart.File) error {
	buff := make([]byte, 512)

	_, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return err
	}

	contentType := http.DetectContentType(buff)

	switch contentType {
	case "image/png", "image/jpeg":
		// allowed
	default:
		return errors.New("Only PNG and JPEG images are allowed.")
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}

func saveUpload(file multipart.File, uploadPath string) error {
	dst, err := os.Create(uploadPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func parseCellSize(r *http.Request, defaultSize int) int {
	cellSize := defaultSize

	if v := r.FormValue("cellsize"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n > 0 {
			cellSize = n
		}
	}

	return cellSize
}

func buildOutputPath(outputDir, uploadedFileName string) string {
	name := strings.TrimSuffix(
		uploadedFileName,
		filepath.Ext(uploadedFileName),
	)

	return filepath.Join(
		outputDir,
		name+"_dot.png",
	)
}

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[32m" // Green
	case code >= 300 && code < 400:
		return "\033[33m" // Yellow
	case code >= 400 && code < 500:
		return "\033[34m" // Blue
	default:
		return "\033[31m" // Red
	}
}

const reset = "\033[0m"
