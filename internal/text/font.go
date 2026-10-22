package text

import (
	"fmt"
	"math"
	"strings"
)

// Style is a combination of Bold and Italic.
type Style uint8

const (
	Bold Style = 1 << iota
	Italic
)

// Font is a resolved metric face at one point size. It uses only built-in
// metric tables and does not depend on installed system fonts.
type Font struct {
	Name   string
	Style  Style
	Size   float64
	family family
}
type family uint8

const (
	sans family = iota
	serif
	mono
)

// Resolve maps logical font names to bundled metric families. Empty and
// unknown names fall back to SansSerif. The returned Name is canonical.
func Resolve(name string, style Style, size float64) (Font, error) {
	if size <= 0 || math.IsNaN(size) || math.IsInf(size, 0) {
		return Font{}, fmt.Errorf("text: invalid font size %v", size)
	}
	if style & ^(Bold|Italic) != 0 {
		return Font{}, fmt.Errorf("text: unsupported font style %d", style)
	}
	f := Font{Style: style, Size: size}
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "serif", "times", "times new roman":
		f.Name = "Serif"
		f.family = serif
	case "monospaced", "monospace", "courier":
		f.Name = "Monospaced"
		f.family = mono
	default:
		f.Name = "SansSerif"
		f.family = sans
	}
	return f, nil
}

// GlyphAdvance returns the nominal advance width in 1/1000 em units. The
// baseline table covers ASCII; other runes use a fixed fallback advance.
func (f Font) GlyphAdvance(r rune) int {
	if f.family == mono {
		return 600
	}
	if r < 32 || r > 126 {
		return 600
	}
	if r == ' ' {
		if f.family == serif {
			return 250
		}
		return 278
	}
	if strings.ContainsRune("ilI.,:;!|'`", r) {
		if f.family == serif {
			return 278
		}
		return 278
	}
	if strings.ContainsRune("MW@", r) {
		if f.family == serif {
			return 889
		}
		return 833
	}
	if r >= '0' && r <= '9' {
		return 556
	}
	if strings.ContainsRune("fjrt()[]{}", r) {
		return 389
	}
	if f.family == serif {
		return 500
	}
	return 556
}
