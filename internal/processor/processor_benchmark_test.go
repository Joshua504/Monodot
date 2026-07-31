package processor

import (
	"path/filepath"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	input := filepath.Join("testdata", "images98.jpeg")

	tests := []struct {
		name     string
		cellSize int
	}{
		{"CellSize2", 2},
		{"CellSize3", 3},
		{"CellSize5", 5},
		{"CellSize8", 8},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			outputDir := b.TempDir()
			output := filepath.Join(outputDir, "output.png")

			b.ResetTimer()

			for b.Loop() {
				if err := Generate(input, output, tt.cellSize); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
