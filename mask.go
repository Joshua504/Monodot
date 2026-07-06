package main

func buildMask(grid [][]float64) [][]bool {
	rows := len(grid)
	cols := len(grid[0])

	mask := make([][]bool, rows)
	for i := range mask {
		mask[i] = make([]bool, cols)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			mask[r][c] = grid[r][c] > 15
		}
	}

	return mask
}
