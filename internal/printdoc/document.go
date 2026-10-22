package printdoc

import "fmt"

// ImagePayload identifies an image resource; resource storage is resolved by
// a later print-document task.
type ImagePayload struct{ Resource string }

// ShapeKind identifies a primitive vector shape.
type ShapeKind uint8

const (
	ShapeLine ShapeKind = iota
	ShapeRectangle
	ShapeEllipse
)

// ShapePayload identifies a vector primitive. Stroke and fill styles are
// supplied by the resolved-style task.
type ShapePayload struct{ Kind ShapeKind }

// NewTextElement creates a positioned text element.
func NewTextElement(bounds Bounds, text string) Element {
	return Element{Kind: ElementText, Bounds: bounds, Text: text}
}

// NewImageElement creates a positioned image reference.
func NewImageElement(bounds Bounds, resource string) Element {
	return Element{Kind: ElementImage, Bounds: bounds, Image: &ImagePayload{Resource: resource}}
}

// NewShapeElement creates a positioned shape.
func NewShapeElement(bounds Bounds, kind ShapeKind) Element {
	return Element{Kind: ElementShape, Bounds: bounds, Shape: &ShapePayload{Kind: kind}}
}

// AppendPage validates and appends a page, assigning a sequential one-based
// number and updating the document's page count. Element slices and payloads
// are copied so the caller can reuse its Page value.
func (d *Document) AppendPage(page Page) error {
	if d == nil {
		return fmt.Errorf("printdoc: nil document")
	}
	page.Number = len(d.Pages) + 1
	if err := page.Validate(); err != nil {
		return err
	}
	for i, elem := range page.Elements {
		if err := validateElement(elem); err != nil {
			return fmt.Errorf("printdoc: page %d element %d: %w", page.Number, i, err)
		}
	}
	copyPage := page
	copyPage.Elements = make([]Element, len(page.Elements))
	for i, elem := range page.Elements {
		copyPage.Elements[i] = elem
		if elem.Image != nil {
			payload := *elem.Image
			copyPage.Elements[i].Image = &payload
		}
		if elem.Shape != nil {
			payload := *elem.Shape
			copyPage.Elements[i].Shape = &payload
		}
	}
	d.Pages = append(d.Pages, copyPage)
	d.Metadata.PageCount = len(d.Pages)
	return nil
}

func validateElement(e Element) error {
	if e.Bounds.Width < 0 || e.Bounds.Height < 0 {
		return fmt.Errorf("negative element size")
	}
	switch e.Kind {
	case ElementText:
		if e.Image != nil || e.Shape != nil {
			return fmt.Errorf("text element has another payload")
		}
	case ElementImage:
		if e.Image == nil || e.Image.Resource == "" || e.Shape != nil {
			return fmt.Errorf("invalid image payload")
		}
	case ElementShape:
		if e.Shape == nil || e.Image != nil || e.Shape.Kind > ShapeEllipse {
			return fmt.Errorf("invalid shape payload")
		}
	default:
		return fmt.Errorf("unknown element kind %d", e.Kind)
	}
	return nil
}

// Validate checks page numbering, metadata, dimensions, and element payloads.
// An empty document with PageCount zero is valid.
func (d *Document) Validate() error {
	if d == nil {
		return fmt.Errorf("printdoc: nil document")
	}
	if d.Metadata.PageCount != len(d.Pages) {
		return fmt.Errorf("printdoc: page count %d differs from %d pages", d.Metadata.PageCount, len(d.Pages))
	}
	for i, page := range d.Pages {
		if page.Number != i+1 {
			return fmt.Errorf("printdoc: page %d has number %d", i+1, page.Number)
		}
		if err := page.Validate(); err != nil {
			return err
		}
		for j, elem := range page.Elements {
			if err := validateElement(elem); err != nil {
				return fmt.Errorf("printdoc: page %d element %d: %w", i+1, j, err)
			}
		}
	}
	return nil
}

// ForEachPage visits pages in document order and stops on the first error.
func (d *Document) ForEachPage(visit func(Page) error) error {
	if d == nil || visit == nil {
		return fmt.Errorf("printdoc: nil document or visitor")
	}
	for _, page := range d.Pages {
		if err := visit(page); err != nil {
			return err
		}
	}
	return nil
}

// ForEachElement visits elements in page order and returns the first error.
func (d *Document) ForEachElement(visit func(int, Element) error) error {
	if d == nil || visit == nil {
		return fmt.Errorf("printdoc: nil document or visitor")
	}
	for _, page := range d.Pages {
		for _, elem := range page.Elements {
			if err := visit(page.Number, elem); err != nil {
				return err
			}
		}
	}
	return nil
}
