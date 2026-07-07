package processor

import (
	"image"
	"image/png"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

func Generate(inputPath, outputPath string, cellsize int) error {

	if cellsize <= 0 {
		cellsize = 3
	}
	openImg, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer openImg.Close()

	img, _, err := image.Decode(openImg)
	if err != nil {
		return err
	}

	grid := buildBrightnessGrid(img, cellsize)
	mask := buildMask(grid)
	canvas := renderDots(grid, mask, cellsize)

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	err = png.Encode(output, canvas)
	if err != nil {
		return err
	}

	return nil
}
