package printdoc

import (
	"fmt"
	"math"
)

// Unit is an internal measurement in twips (1/1440 inch).
//
// The twip was chosen as the print IR unit because it is (a) integral, so
// repeated layout arithmetic does not accumulate floating-point error, and (b)
// exactly divisible into both printers' points (1 pt = 20 twips) and inches
// (1 in = 1440 twips). See docs/adr/0004-print-ir-unit-system.md.
type Unit int64

const (
	// TwipsPerInch is the number of twips in one inch.
	TwipsPerInch Unit = 1440
	// TwipsPerPoint is the number of twips in one point (1/72 inch).
	TwipsPerPoint Unit = 20
)

// FromPoints converts points to the internal unit, rounding half away from
// zero.
func FromPoints(points float64) Unit { return Unit(math.Round(points * float64(TwipsPerPoint))) }

// FromInches converts inches to the internal unit, rounding half away from
// zero.
func FromInches(inches float64) Unit { return Unit(math.Round(inches * float64(TwipsPerInch))) }

// FromMM converts millimetres to the internal unit, rounding half away from
// zero.
func FromMM(mm float64) Unit { return Unit(math.Round(mm * float64(TwipsPerInch) / 25.4)) }

// Points returns the value in points.
func (u Unit) Points() float64 { return float64(u) / float64(TwipsPerPoint) }

// Inches returns the value in inches.
func (u Unit) Inches() float64 { return float64(u) / float64(TwipsPerInch) }

// MM returns the value in millimetres.
func (u Unit) MM() float64 { return float64(u) * 25.4 / float64(TwipsPerInch) }

// Bounds is an axis-aligned rectangle in page coordinates.
type Bounds struct {
	X      Unit
	Y      Unit
	Width  Unit
	Height Unit
}

// ElementKind classifies a print element.
type ElementKind int

const (
	// ElementText is a run of text.
	ElementText ElementKind = iota
	// ElementImage is a raster image reference.
	ElementImage
	// ElementShape is a vector shape.
	ElementShape
)

// String returns the lowercase element kind name.
func (k ElementKind) String() string {
	switch k {
	case ElementText:
		return "text"
	case ElementImage:
		return "image"
	case ElementShape:
		return "shape"
	default:
		return fmt.Sprintf("elementKind(%d)", int(k))
	}
}

// Element is a positioned element on a page. Only the fields relevant to Kind
// are meaningful; richer payloads (fonts, images, shapes) are added by later
// tasks (PRT-002, PRT-003).
type Element struct {
	Kind   ElementKind
	Bounds Bounds
	// Text is set when Kind is ElementText.
	Text string
	// Image is set when Kind is ElementImage.
	Image *ImagePayload
	// Shape is set when Kind is ElementShape.
	Shape *ShapePayload
}

// Page is a single printed page.
type Page struct {
	Number   int
	Width    Unit
	Height   Unit
	Elements []Element
}

// Validate reports whether the page has a positive size.
func (p Page) Validate() error {
	if p.Width <= 0 || p.Height <= 0 {
		return fmt.Errorf("printdoc: page %d has non-positive size %dx%d", p.Number, p.Width, p.Height)
	}
	return nil
}

// Metadata describes a whole document.
type Metadata struct {
	// PageCount is the number of pages in the document. Per-page numbers live
	// on Page.Number.
	PageCount int
}

// Document is a paginated print document. Pages are ordered; Metadata.PageCount
// should equal len(Pages).
type Document struct {
	Pages    []Page
	Metadata Metadata
}
