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
	defer func() {
		if err := openImg.Close(); err != nil {
			log.Printf("failed to close input image: %v", err)
		}
	}()

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
	defer func() {
		if err := output.Close(); err != nil {
			log.Printf("failed to close output image: %v", err)
		}
	}()

	err = png.Encode(output, canvas)
	if err != nil {
		return err
	}

	return nil
}
