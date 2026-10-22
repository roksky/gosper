package text

import (
	"fmt"
	"github.com/roksky/gosper/internal/printdoc"
	"math"
	"unicode/utf8"
)

const maxMeasureBytes = 1 << 20

// Metrics holds baseline measurements in print-document twips.
type Metrics struct {
	Width      printdoc.Unit
	Ascent     printdoc.Unit
	Descent    printdoc.Unit
	LineHeight printdoc.Unit
}

// Measure returns deterministic, single-line metrics for s using f. It sums
// nominal rune advances and rounds once to twips; it does not shape, kern,
// wrap, or cache text. Empty text has zero width and normal vertical metrics.
func Measure(s string, f Font) (Metrics, error) {
	if f.Size <= 0 || f.Size > 10000 || math.IsNaN(f.Size) || math.IsInf(f.Size, 0) {
		return Metrics{}, fmt.Errorf("text: invalid font size %v", f.Size)
	}
	if len(s) > maxMeasureBytes {
		return Metrics{}, fmt.Errorf("text: text exceeds %d bytes", maxMeasureBytes)
	}
	if !utf8.ValidString(s) {
		return Metrics{}, fmt.Errorf("text: invalid UTF-8")
	}
	var advance int64
	for _, r := range s {
		if r == '\n' || r == '\r' {
			return Metrics{}, fmt.Errorf("text: multiline input requires line layout")
		}
		advance += int64(f.GlyphAdvance(r))
	}
	twips := float64(printdoc.TwipsPerPoint)
	return Metrics{
		Width:      printdoc.Unit(math.Round(float64(advance) * f.Size * twips / 1000)),
		Ascent:     printdoc.Unit(math.Round(.8 * f.Size * twips)),
		Descent:    printdoc.Unit(math.Round(.2 * f.Size * twips)),
		LineHeight: printdoc.Unit(math.Round(1.2 * f.Size * twips)),
	}, nil
}
