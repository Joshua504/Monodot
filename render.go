package main

import (
	"image"
	"image/color"
	"image/draw"
)

func renderDots(grid [][]float64, cellSize int) *image.RGBA {
	row := len(grid)
	col := len(grid[0])
	width := col * cellSize
	height := row * cellSize

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	for numRow := 0; numRow < row; numRow++ {
		for numCol := 0; numCol < col; numCol++ {
			if shouldPlaceDot(numRow, numCol, grid[numRow][numCol]) {
				cx := numCol*cellSize + cellSize/2
				cy := numRow*cellSize + cellSize/2
				drawCircle(canvas, cx, cy, 1, color.Black)
			}
		}
	}

	return canvas
}

func drawCircle(canvas *image.RGBA, cx, cy, radius int, col color.Color) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= radius*radius {
				canvas.Set(x, y, col)
			}
		}
	}
}
