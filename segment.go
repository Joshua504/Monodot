package main

import "image"

func averageSaturation(img image.Image, x0, y0, cellSize int) float64 {
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

			max := r8
			if g8 > max {
				max = g8
			}
			if b8 > max {
				max = b8
			}

			min := r8
			if g8 < min {
				min = g8
			}
			if b8 < min {
				min = b8
			}

			var sat float64
			if max > 0 {
				sat = (max - min) / max
			} else {
				sat = 0
			}

			total += sat
			count++
		}
	}
	if count == 0 {
		return 0
	}

	return total / float64(count)
}
