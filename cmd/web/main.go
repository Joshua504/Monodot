package main

import (
	"net/http"
	"path/filepath"
)

func main() {
	cfg := NewConfig()

	startCleanupJob()

	server := newServer(cfg)

	go startServer(server)

	waitForShutdown(server)
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	app.render(
		w,
		http.StatusOK,
		"index.html",
		nil,
	)
}

func (app *application) generateHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info(
		"Upload request received",
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
	)

	if r.Method != http.MethodPost {
		app.logger.Warn(
			"Invalid request method",
			"method", r.Method,
		)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(app.config.MaxUploadSize)
	if err != nil {
		app.logger.Warn(
			"Failed to parse multipart form",
			"error", err,
		)
		app.validationError(w, "Invalid multipart for")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		app.validationError(w, "No image was uploaded.")
		return
	}
	app.logger.Info(
		"Upload received",
		"filename", header.Filename,
	)
	defer file.Close()

	err = validateExtension(header.Filename)
	if err != nil {
		app.logger.Warn(
			"Extension validation failed",
			"filename", header.Filename,
			"error", err,
		)
		app.validationError(w, err.Error())
		return
	}

	err = validateContentType(file)
	if err != nil {
		app.logger.Warn(
			"MIME validation failed",
			"filename", header.Filename,
			"error", err,
		)
		app.validationError(w, err.Error())
		return
	}

	uploadedFileName := uniqueName(header.Filename)
	uploadPath := filepath.Join(app.config.UploadDir, uploadedFileName)

	outputPath := buildOutputPath(app.config.OutputDir, uploadedFileName)

	err = saveUpload(file, uploadPath)
	app.logger.Info(
		"Upload saved",
		"path", uploadPath,
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	cellsize := parseCellSize(r, app.config.DefaultCellSize)

	err = app.imageGenerator(
		uploadPath,
		outputPath,
		cellsize,
	)
	app.logger.Info(
		"Starting image generation",
		"cell_size", cellsize,
	)
	if err != nil {
		app.logger.Error(
			"Image generation failed",
			"error", err,
		)
		app.serverError(w, err)
		return
	}
	app.logger.Info(
		"Image generated successfully",
		"output", outputPath,
	)

	app.logger.Info(
		"Redirecting to result page",
		"image", filepath.Base(outputPath),
	)
	http.Redirect(
		w,
		r,
		"/result?image=/outputs/"+filepath.Base(outputPath),
		http.StatusSeeOther,
	)
}

func (app *application) resultHandler(w http.ResponseWriter, r *http.Request) {
	image := r.URL.Query().Get("image")

	data := map[string]string{
		"Image": image,
	}

	app.render(
		w,
		http.StatusOK,
		"result.html",
		data,
	)
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		app.logger.Error(
			"Failed to write health response",
			"error", err,
		)
	}
}
