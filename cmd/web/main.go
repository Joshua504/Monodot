package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Joshua504/Monodot/internal/processor"
)

var tmpl = template.Must(
	template.ParseGlob("templates/*.html"))

func main() {
	cfg := NewConfig()

	startCleanupJob()

	server := newServer(cfg)

	go startServer(server)

	waitForShutdown(server)
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) generateHandler(w http.ResponseWriter, r *http.Request) {
	logger := requestLogger(r)
	logger.Println("Upload request received")

	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Printf("Failed to parse multipart form: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	logger.Printf("Upload received: %s", header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = validateExtension(header.Filename)
	if err != nil {
		logger.Printf("Extension validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateContentType(file)
	if err != nil {
		logger.Printf("MIME validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadedFileName := uniqueName(header.Filename)
	uploadPath := filepath.Join(app.config.UploadDir, uploadedFileName)

	outputPath := buildOutputPath(app.config.OutputDir, uploadedFileName)

	err = saveUpload(file, uploadPath)
	logger.Printf("Upload saved: %s", uploadPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cellsize := parseCellSize(r, app.config.DefaultCellSize)

	err = processor.Generate(uploadPath, outputPath, cellsize)
	log.Printf("Starting image generation")
	if err != nil {
		logger.Printf("Image generation failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Printf("Image generated successfully: %s", outputPath)

	logger.Printf("Redirecting to result page")
	http.Redirect(
		w,
		r,
		"/result?image=/outputs/"+filepath.Base(outputPath),
		http.StatusSeeOther,
	)
}

func (app *application) resultHandler(w http.ResponseWriter, r *http.Request) {
	image := r.URL.Query().Get("image")

	err := tmpl.ExecuteTemplate(w, "result.html", struct{ Image string }{Image: image})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
