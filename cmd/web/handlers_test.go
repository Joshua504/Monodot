package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	cfg := NewConfig()

	cache, err := newTemplateCache(cfg.TemplateDir)
	if err != nil {
		t.Fatal(err)
	}

	logger := log.New(io.Discard, "", 0)

	return &application{
		config:        cfg,
		logger:        logger,
		templateCache: cache,
		imageGenerator: func(input, output string, cellSize int) error {
			return nil
		},
	}
}

func TestHealthHandler(t *testing.T) {
	app := newTestApplication(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	rr := httptest.NewRecorder()
	app.healthHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d",
			http.StatusOK,
			res.StatusCode)
	}

	want := "application/json"
	if got := res.Header.Get("Content-Type"); got != want {
		t.Errorf("expected Content-Type %q, got %q",
			want,
			got)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	wantBody := `{"status":"ok"}`

	if string(body) != wantBody {
		t.Errorf(
			"expected body %q, got %q",
			wantBody,
			string(body),
		)
	}
}

func TestHomeHandler(t *testing.T) {
	app := newTestApplication(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rr := httptest.NewRecorder()

	app.homeHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(body) == 0 {
		t.Error("expected response body to contain rendered HTML")
	}

	if !strings.Contains(string(body), "<html") {
		t.Error("expected rendered HTML document")
	}
}

func TestResultHandler(t *testing.T) {
	app := newTestApplication(t)

	imagePath := "/outputs/test_image_dot.png"

	req := httptest.NewRequest(
		http.MethodGet,
		"/result?image="+url.QueryEscape(imagePath),
		nil,
	)

	rr := httptest.NewRecorder()

	app.resultHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if len(html) == 0 {
		t.Fatal("expected HTML response")
	}

	if !strings.Contains(html, imagePath) {
		t.Errorf(
			"expected rendered HTML to contain %q",
			imagePath,
		)
	}
}

func TestGenerateHandlerInvalidMethod(t *testing.T) {
	app := newTestApplication(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/generate",
		nil,
	)

	rr := httptest.NewRecorder()

	app.generateHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			res.StatusCode,
		)
	}
}

func TestGenerateHandlerMissingImage(t *testing.T) {
	app := newTestApplication(t)

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)
	writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/generate",
		&body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	rr := httptest.NewRecorder()

	app.generateHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}
}

func TestGenerateHandlerInvalidExtension(t *testing.T) {
	app := newTestApplication(t)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", "hello.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = part.Write([]byte("this is not an image"))
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/generate",
		&body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	rr := httptest.NewRecorder()

	app.generateHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}
}

func TestGenerateHandlerInvalidMIMEType(t *testing.T) {
	app := newTestApplication(t)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", "image.png")
	if err != nil {
		t.Fatal(err)
	}

	// This is NOT a PNG.
	_, err = part.Write([]byte("I am definitely not an image"))
	if err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/generate",
		&body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	rr := httptest.NewRecorder()

	app.generateHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}
}

func createTestPNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	var buf bytes.Buffer

	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}

	return buf.Bytes()
}

func TestGenerateHandlerSuccess(t *testing.T) {
	cfg := NewConfig()

	cache, err := newTemplateCache(cfg.TemplateDir)
	if err != nil {
		t.Fatal(err)
	}

	called := false

	app := &application{
		config:        cfg,
		logger:        log.New(os.Stdout, "", log.LstdFlags),
		templateCache: cache,
		imageGenerator: func(input, output string, cellSize int) error {
			called = true
			return nil
		},
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", "test.png")
	if err != nil {
		t.Fatal(err)
	}

	_, err = part.Write(createTestPNG(t))
	if err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/generate",
		&body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	rr := httptest.NewRecorder()

	app.generateHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusSeeOther,
			res.StatusCode,
		)
	}

	location := res.Header.Get("Location")

	if !strings.HasPrefix(location, "/result?image=/outputs/") {
		t.Fatalf(
			"unexpected redirect location: %s",
			location,
		)
	}

	_ = called
}
