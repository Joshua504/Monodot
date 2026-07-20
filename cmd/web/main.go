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
	app.logger.Println("Upload request received")

	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(app.config.MaxUploadSize)
	if err != nil {
		app.logger.Printf("Failed to parse multipart form: %v", err)
		app.validationError(w, "Invalid multipart for")
		return
	}

	file, header, err := r.FormFile("image")
	app.logger.Printf("Upload received: %s", header.Filename)
	if err != nil {
		app.validationError(w, "No image was uploaded.")
		return
	}
	defer file.Close()

	err = validateExtension(header.Filename)
	if err != nil {
		app.logger.Printf("Extension validation failed: %v", err)
		app.validationError(w, err.Error())
		return
	}

	err = validateContentType(file)
	if err != nil {
		app.logger.Printf("MIME validation failed: %v", err)
		app.validationError(w, err.Error())
		return
	}

	uploadedFileName := uniqueName(header.Filename)
	uploadPath := filepath.Join(app.config.UploadDir, uploadedFileName)

	outputPath := buildOutputPath(app.config.OutputDir, uploadedFileName)

	err = saveUpload(file, uploadPath)
	app.logger.Printf("Upload saved: %s", uploadPath)
	if err != nil {
		app.serverError(w, err)
		return
	}

	cellsize := parseCellSize(r, app.config.DefaultCellSize)

	err = processor.Generate(uploadPath, outputPath, cellsize)
	log.Printf("Starting image generation")
	if err != nil {
		app.logger.Printf("Image generation failed: %v", err)
		app.serverError(w, err)
		return
	}
	app.logger.Printf("Image generated successfully: %s", outputPath)

	app.logger.Printf("Redirecting to result page")
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
