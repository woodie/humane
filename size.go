// Package humane formats file sizes and relative dates the way macOS Finder does.
package humane

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// sizeUnits are Finder's capitalized, 1000-based unit labels above byte scale.
var sizeUnits = []string{"KB", "MB", "GB", "TB", "PB", "EB"}

// HumanSize formats a byte count the way Finder does. See docs/COMMENTS.md
// for the 3-significant-figure rounding rule and the corrections (zero
// bytes, byte-scale wording) this bakes in from real ByteCountFormatter
// output.
func HumanSize(bytes int64) string {
	if bytes == 0 {
		return "Zero KB"
	}
	if bytes < 1000 {
		if bytes == 1 {
			return "1 byte"
		}
		return fmt.Sprintf("%d bytes", bytes)
	}

	exp := int(math.Log(float64(bytes)) / math.Log(1000))
	if exp > len(sizeUnits) {
		exp = len(sizeUnits)
	}
	value := float64(bytes) / math.Pow(1000, float64(exp))

	return fmt.Sprintf("%s %s", formatSignificant(value, 3), sizeUnits[exp-1])
}

// formatSignificant rounds value to sigFigs significant figures and trims
// any trailing fractional zeros (and the decimal point itself, if nothing
// remains after it) -- see docs/COMMENTS.md.
func formatSignificant(value float64, sigFigs int) string {
	magnitude := int(math.Floor(math.Log10(value))) + 1
	decimals := sigFigs - magnitude
	if decimals < 0 {
		decimals = 0
	}

	s := strconv.FormatFloat(value, 'f', decimals, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
