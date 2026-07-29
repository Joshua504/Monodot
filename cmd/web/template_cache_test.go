package main

import "testing"

func TestNewTemplateCache(t *testing.T) {
	cache, err := newTemplateCache("templates")
	if err != nil {
		t.Fatalf("newTemplateCache() returned error: %v", err)
	}

	if cache == nil {
		t.Fatal("expected cache to be initialized")
	}

	if len(cache) == 0 {
		t.Fatal("expected cache to contain templates")
	}
}

func TestTemplateCacheContainsPages(t *testing.T) {
	cache, err := newTemplateCache("templates")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"index.html",
		"result.html",
		"404.html",
		"500.html",
	}

	for _, page := range expected {
		if _, ok := cache[page]; !ok {
			t.Errorf("template %q not found in cache", page)
		}
	}

}
