package utils

import (
	"strings"
)

func Crop(s string, n int) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if len(l) >= 2 {
			lines[i] = l[n:]
		}
	}
	return strings.Join(lines, "\n")
}
