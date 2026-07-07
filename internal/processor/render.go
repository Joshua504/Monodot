package processor

import (
	"image"
	"image/color"
	"image/draw"
)

func renderDots(grid [][]float64, mask [][]bool, cellSize int) *image.RGBA {
	row := len(grid)
	col := len(grid[0])
	width := col * cellSize
	height := row * cellSize

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	for numRow := 0; numRow < row; numRow++ {
		for numCol := 0; numCol < col; numCol++ {
			brightness := grid[numRow][numCol]

			if !isForeground(brightness) {
				continue
			}
			if !mask[numRow][numCol] {
				continue
			}
			if brightness > 245 {
				continue
			}

			cx := numCol*cellSize + cellSize/2
			if numRow%2 == 1 {
				cx += cellSize / 2
			}
			cy := numRow*cellSize + cellSize/2

			radius := 1

			gray := uint8(brightness)

			dotColor := color.RGBA{gray, gray, gray, 255}

			drawCircle(canvas, cx, cy, radius, dotColor)
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

func isForeground(brightness float64) bool {
	return brightness > 3
}
