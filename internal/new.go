package new

import (
	"math"
	"strings"
	"unicode"
)

func Encode(pt string) string {

	var normalized []rune

	for _, ch := range pt {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			normalized = append(normalized, unicode.ToLower(ch))
		}
	}

	if len(normalized) == 0 {
		return ""
	}

	length := len(normalized)

	c := int(math.Ceil(math.Sqrt(float64(length))))
	r := int(math.Ceil(float64(length) / float64(c)))

	grid := make([][]rune, r)

	for i := 0; i < r; i++ {
		grid[i] = normalized[i*c : (i+1)*c]
	}

	chunks := make([]string, c)
	for col := 0; col < c; col++ {
		col_chars := make([]rune, c)
		for row := 0; row < r; row++ {
			col_chars[row] = grid[row][col]
		}
		chunks[col] = string(col_chars)
	}

	return strings.Join(chunks, "")
}
