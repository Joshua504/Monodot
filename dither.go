package main

var bayerMatrix4x4 = [4][4]int{
	{0, 8, 2, 10},
	{12, 4, 14, 6},
	{3, 11, 1, 9},
	{15, 7, 5, 13},
}

func shouldPlaceDot(row, col int, brightness float64) bool {
	threshold := float64(bayerMatrix4x4[row%4][col%4]) * (255.0 / 16.0)

	if brightness < threshold {
		return true
	}

	return false
}
