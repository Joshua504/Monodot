package main

type ResultPageData struct {
	Image    string
	Original string
	FileName string

	Width  int
	Height int

	OriginalSize  string
	GeneratedSize string

	CellSize int

	ProcessingTime string
}
