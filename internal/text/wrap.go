package text

import (
	"fmt"
	"strings"

	"github.com/roksky/gosper/internal/printdoc"
)

// Wrap greedily fits words within width. Whitespace between words collapses to
// one space; explicit newlines start a new line. An oversized word is split at
// rune boundaries. A glyph wider than width is an error.
func Wrap(s string, font Font, width printdoc.Unit) ([]string, error) {
	if width < 0 {
		return nil, fmt.Errorf("text: negative wrap width")
	}
	if len(s) > maxMeasureBytes {
		return nil, fmt.Errorf("text: text exceeds %d bytes", maxMeasureBytes)
	}
	if s == "" {
		return []string{""}, nil
	}
	var lines []string
	for _, paragraph := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}
		line := ""
		for _, word := range words {
			candidate := word
			if line != "" {
				candidate = line + " " + word
			}
			m, err := Measure(candidate, font)
			if err != nil {
				return nil, err
			}
			if m.Width <= width {
				line = candidate
				continue
			}
			if line != "" {
				lines = append(lines, line)
				line = ""
			}
			m, err = Measure(word, font)
			if err != nil {
				return nil, err
			}
			if m.Width <= width {
				line = word
				continue
			}
			for _, glyph := range word {
				candidate = line + string(glyph)
				m, err := Measure(candidate, font)
				if err != nil {
					return nil, err
				}
				if m.Width > width {
					if line == "" {
						return nil, fmt.Errorf("text: glyph exceeds wrap width")
					}
					lines = append(lines, line)
					line = string(glyph)
					if m, err = Measure(line, font); err != nil || m.Width > width {
						return nil, fmt.Errorf("text: glyph exceeds wrap width")
					}
				} else {
					line = candidate
				}
			}
		}
		lines = append(lines, line)
	}
	return lines, nil
}
