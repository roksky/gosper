package printdoc

import "fmt"

// Align is horizontal text alignment.
type Align uint8

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Insets are element padding in twips.
type Insets struct{ Top, Right, Bottom, Left Unit }

// Style contains resolved visual values carried by a print element. Color is
// an opaque 0xRRGGBB RGB value; no styles or palette names reach exporters.
type Style struct {
	Font         string
	Size         Unit
	Bold, Italic bool
	Color        uint32
	Align        Align
	Padding      Insets
}

// DefaultStyle is the M1 text style: 10-point SansSerif, black, left-aligned,
// without padding.
func DefaultStyle() Style {
	return Style{Font: "SansSerif", Size: FromPoints(10), Color: 0, Align: AlignLeft}
}

// EffectiveStyle substitutes defaults for an entirely absent style, and fills
// only omitted font and size on a partially specified style.
func (e Element) EffectiveStyle() Style {
	s := e.Style
	if s == (Style{}) {
		return DefaultStyle()
	}
	if s.Font == "" {
		s.Font = DefaultStyle().Font
	}
	if s.Size == 0 {
		s.Size = DefaultStyle().Size
	}
	return s
}

// Validate checks the resolved M1 subset. Canonical font names are required.
func (s Style) Validate() error {
	if s.Size <= 0 || s.Size > FromPoints(10000) {
		return fmt.Errorf("printdoc: invalid style size %d", s.Size)
	}
	switch s.Font {
	case "SansSerif", "Serif", "Monospaced":
	default:
		return fmt.Errorf("printdoc: unresolved font %q", s.Font)
	}
	if s.Color > 0xFFFFFF {
		return fmt.Errorf("printdoc: invalid RGB color %#x", s.Color)
	}
	if s.Align > AlignRight {
		return fmt.Errorf("printdoc: invalid alignment %d", s.Align)
	}
	if s.Padding.Top < 0 || s.Padding.Right < 0 || s.Padding.Bottom < 0 || s.Padding.Left < 0 {
		return fmt.Errorf("printdoc: negative padding")
	}
	return nil
}
