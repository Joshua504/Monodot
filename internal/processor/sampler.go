package processor

import "image"

func averageBrightness(img image.Image, x0, y0, cellSize int) float64 {
	bounds := img.Bounds()
	xEnd := x0 + cellSize
	if xEnd > bounds.Max.X {
		xEnd = bounds.Max.X
	}

	yEnd := y0 + cellSize
	if yEnd > bounds.Max.Y {
		yEnd = bounds.Max.Y
	}

	var total float64
	var count int

	for y := y0; y < yEnd; y++ {
		for x := x0; x < xEnd; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			gray := 0.299*r8 + 0.587*g8 + 0.114*b8
			total += gray
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func buildBrightnessGrid(img image.Image, cellSize int) [][]float64 {
	bounds := img.Bounds()

	cols := (bounds.Dx() + cellSize - 1) / cellSize
	rows := (bounds.Dy() + cellSize - 1) / cellSize

	grid := make([][]float64, rows)

	for i := range grid {
		grid[i] = make([]float64, cols)
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			x0 := col * cellSize
			y0 := row * cellSize
			grid[row][col] = averageBrightness(img, x0, y0, cellSize)
		}
	}

	return grid
}
