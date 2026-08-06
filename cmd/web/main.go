package main

import (
	"image"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	cfg := NewConfig()

	startCleanupJob()

	server := newServer(cfg)

	go startServer(server)

	waitForShutdown(server)
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

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
	defer func() {
		if err := file.Close(); err != nil {
			app.logger.Error(
				"failed to close uploaded file",
				"error", err,
			)
		}
	}()

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

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	width := config.Width
	height := config.Height

	_, err = file.Seek(0, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

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

	start := time.Now()

	err = app.imageGenerator(
		uploadPath,
		outputPath,
		cellsize,
	)

	processingTime := time.Since(start)

	uploadInfo, err := os.Stat(uploadPath)
	if err != nil {
		app.serverError(w, err)
		return
	}

	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		app.serverError(w, err)
		return
	}

	originalSize := formatFileSize(uploadInfo.Size())
	generatedSize := formatFileSize(outputInfo.Size())

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
	params := url.Values{}

	params.Set(
		"image",
		"/outputs/"+filepath.Base(outputPath),
	)

	params.Set(
		"original",
		"/uploads/"+uploadedFileName,
	)

	params.Set(
		"name",
		header.Filename,
	)

	params.Set(
		"cellsize",
		r.FormValue("cellSize"),
	)

	params.Set(
		"time",
		processingTime.Round(time.Millisecond).String(),
	)

	params.Set(
		"width",
		strconv.Itoa(width),
	)

	params.Set(
		"height",
		strconv.Itoa(height),
	)

	params.Set(
		"originalsize",
		originalSize,
	)

	params.Set(
		"generatedsize",
		generatedSize,
	)

	http.Redirect(
		w,
		r,
		"/result?"+params.Encode(),
		http.StatusSeeOther,
	)
}

func (app *application) resultHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	image := r.URL.Query().Get("image")
	original := r.URL.Query().Get("original")
	fileName := r.URL.Query().Get("name")
	processingTime := r.URL.Query().Get("time")

	cellSize, _ := strconv.Atoi(
		r.URL.Query().Get("cellsize"),
	)

	width, _ := strconv.Atoi(
		r.URL.Query().Get("width"),
	)

	height, _ := strconv.Atoi(
		r.URL.Query().Get("height"),
	)

	originalSize := r.URL.Query().Get("originalsize")

	generatedSize := r.URL.Query().Get("generatedsize")

	data := ResultPageData{
		Image:    image,
		Original: original,
		FileName: fileName,

		Width:  width,
		Height: height,

		OriginalSize:  originalSize,
		GeneratedSize: generatedSize,

		CellSize:       cellSize,
		ProcessingTime: processingTime,
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
