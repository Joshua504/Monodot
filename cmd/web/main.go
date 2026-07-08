package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Joshua504/Monodot/internal/processor"
)

var tmpl = template.Must(
	template.ParseGlob("templates/*.html"))

func main() {
	startCleanupJob()

	server := newServer()

	go startServer(server)

	waitForShutdown(server)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = validateExtension(header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateContentType(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadedFileName := uniqueName(header.Filename)
	uploadPath := filepath.Join("uploads", uploadedFileName)

	outputPath := buildOutputPath(uploadedFileName)

	err = saveUpload(file, uploadPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cellsize := parseCellSize(r)

	err = processor.Generate(uploadPath, outputPath, cellsize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/result?image=/outputs/"+filepath.Base(outputPath),
		http.StatusSeeOther,
	)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	image := r.URL.Query().Get("image")

	err := tmpl.ExecuteTemplate(w, "result.html", struct{ Image string }{Image: image})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
