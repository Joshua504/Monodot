package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestValidateExtension(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "Valid JPG",
			fileName: "image.jpg",
			wantErr:  false,
		},
		{
			name:     "Valid JPEG",
			fileName: "image.jpg",
			wantErr:  false,
		},
		{
			name:     "Valid PNG",
			fileName: "image.jpg",
			wantErr:  false,
		},
		{
			name:     "Uppercase Extension",
			fileName: "image.jpg",
			wantErr:  false,
		},
		{
			name:     "Unsupported GIF",
			fileName: "image.GIF",
			wantErr:  true,
		},
		{
			name:     "Unsupported TXT",
			fileName: "image",
			wantErr:  true,
		},
		{
			name:     "No Extension",
			fileName: "image",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExtension(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExtension(%q) error = %v, wantErr %v", tt.fileName, err, tt.wantErr)
			}
		})
	}
}

func TestParseCellSize(t *testing.T) {
	tests := []struct {
		name        string
		formValue   string
		defaultSize int
		want        int
	}{
		{
			name:        "Valid value",
			formValue:   "8",
			defaultSize: 3,
			want:        8,
		},
		{
			name:        "Empty value uses default",
			formValue:   "",
			defaultSize: 3,
			want:        3,
		},
		{
			name:        "Negative value uses default",
			formValue:   "-5",
			defaultSize: 3,
			want:        3,
		},
		{
			name:        "Zero uses default",
			formValue:   "0",
			defaultSize: 3,
			want:        3,
		},
		{
			name:        "Invalid number uses default",
			formValue:   "abc",
			defaultSize: 3,
			want:        3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)

			q := req.URL.Query()
			q.Add("cellsize", tt.formValue)
			req.URL.RawQuery = q.Encode()

			got := parseCellSize(req, tt.defaultSize)
			if got != tt.want {
				t.Errorf("parseCellSize() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBuildOutputPath(t *testing.T) {
	tests := []struct {
		name             string
		outputDir        string
		uploadedFileName string
		want             string
	}{
		{
			name:             "JPEG file",
			outputDir:        "outputs",
			uploadedFileName: "photo.jpeg",
			want: filepath.Join(
				"outputs",
				"photo_dot.png",
			),
		},
		{
			name:             "PNG file",
			outputDir:        "outputs",
			uploadedFileName: "image.png",
			want: filepath.Join(
				"outputs",
				"image_dot.png",
			),
		},
		{
			name:             "Unique filename",
			outputDir:        "outputs",
			uploadedFileName: "123456_cat.jpg",
			want: filepath.Join(
				"outputs",
				"123456_cat_dot.png",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildOutputPath(tt.outputDir, tt.uploadedFileName)

			if got != tt.want {
				t.Errorf("buildOutputPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUniqueName(t *testing.T) {

}
