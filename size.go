// Package humane formats file sizes and relative dates the way macOS Finder does.
package humane

import (
	"fmt"
	"math"
)

// sizeUnits are Finder's capitalized, 1000-based unit labels.
var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// SizeFormatter formats byte counts the way Finder does; the zero value is ready to use.
type SizeFormatter struct{}

// Format returns bytes as a Finder-style human-readable string.
func (f SizeFormatter) Format(bytes int64) string {
	if bytes < 1000 {
		return fmt.Sprintf("%d B", bytes)
	}
	exp := int(math.Log(float64(bytes)) / math.Log(1000))
	if exp >= len(sizeUnits) {
		exp = len(sizeUnits) - 1
	}
	rounded := math.Round(float64(bytes)/math.Pow(1000, float64(exp))*10) / 10
	if rounded < 10 {
		return fmt.Sprintf("%.1f %s", rounded, sizeUnits[exp])
	}
	return fmt.Sprintf("%.0f %s", rounded, sizeUnits[exp])
}
