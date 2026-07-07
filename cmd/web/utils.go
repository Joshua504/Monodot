package main

import (
	"fmt"
	"time"
)

func uniqueName(filename string) string {
	return fmt.Sprintf("%d_%s",
		time.Now().UnixNano(),
		filename)
}
