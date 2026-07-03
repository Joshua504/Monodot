package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

func main() {
	file := "assets/images.jpeg"

	openImg, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer openImg.Close()

	img, format, err := image.Decode(openImg)
	if err != nil {
		log.Fatalf("Failed to Decode image: %v", err)
	}

	grid := buildBrightnessGrid(img, 4)
	canvas := renderDots(grid, 4)

	output, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer output.Close()
	png.Encode(output, canvas)

	fmt.Println(averageSaturation(img, 250, 150, 8))
	fmt.Println(averageSaturation(img, 0, 0, 8))
	fmt.Println(format)
}
