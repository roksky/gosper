package printdoc

import (
	"errors"
	"testing"
)

func TestDocumentAppendAndIterate(t *testing.T) {
	var d Document
	b := Bounds{X: 10, Y: 20, Width: 30, Height: 40}
	page := Page{Width: FromPoints(100), Height: FromPoints(200), Elements: []Element{NewTextElement(b, "hello"), NewImageElement(b, "logo"), NewShapeElement(b, ShapeRectangle)}}
	if err := d.AppendPage(page); err != nil {
		t.Fatal(err)
	}
	page.Elements[0].Text = "changed"
	page.Elements[1].Image.Resource = "other"
	if err := d.AppendPage(Page{Width: FromPoints(200), Height: FromPoints(300)}); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal(err)
	}
	if d.Metadata.PageCount != 2 || d.Pages[0].Number != 1 || d.Pages[1].Number != 2 || d.Pages[1].Width != FromPoints(200) {
		t.Fatalf("pages=%+v metadata=%+v", d.Pages, d.Metadata)
	}
	if d.Pages[0].Elements[0].Text != "hello" || d.Pages[0].Elements[1].Image.Resource != "logo" || d.Pages[0].Elements[0].Bounds != b {
		t.Fatalf("elements=%+v", d.Pages[0].Elements)
	}
	var pageCount, elemCount int
	if err := d.ForEachPage(func(p Page) error { pageCount++; return nil }); err != nil {
		t.Fatal(err)
	}
	if err := d.ForEachElement(func(number int, e Element) error {
		if number != 1 {
			t.Errorf("element on page %d", number)
		}
		elemCount++
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if pageCount != 2 || elemCount != 3 {
		t.Fatalf("iterations %d,%d", pageCount, elemCount)
	}
}
func TestEmptyDocumentAndValidation(t *testing.T) {
	var d Document
	if err := d.Validate(); err != nil {
		t.Fatal(err)
	}
	n := 0
	if err := d.ForEachPage(func(Page) error { n++; return nil }); err != nil || n != 0 {
		t.Fatalf("empty iteration=%d,%v", n, err)
	}
	if err := d.AppendPage(Page{Width: 0, Height: 20}); err == nil {
		t.Fatal("invalid page accepted")
	}
	if d.Metadata.PageCount != 0 {
		t.Fatal("invalid page changed count")
	}
	d.Metadata.PageCount = 1
	if err := d.Validate(); err == nil {
		t.Fatal("mismatched count accepted")
	}
}
func TestDocumentRejectsInvalidElementsAndStopsOnError(t *testing.T) {
	var d Document
	bad := Page{Width: 20, Height: 20, Elements: []Element{NewImageElement(Bounds{Width: 1, Height: 1}, "")}}
	if err := d.AppendPage(bad); err == nil {
		t.Fatal("empty image resource accepted")
	}
	bad.Elements = []Element{NewShapeElement(Bounds{Width: 1, Height: 1}, ShapeKind(9))}
	if err := d.AppendPage(bad); err == nil {
		t.Fatal("unknown shape accepted")
	}
	bad.Elements = []Element{NewTextElement(Bounds{Width: -1}, "x")}
	if err := d.AppendPage(bad); err == nil {
		t.Fatal("negative bounds accepted")
	}
	_ = d.AppendPage(Page{Width: 20, Height: 20, Elements: []Element{NewTextElement(Bounds{}, "x")}})
	boom := errors.New("stop")
	if err := d.ForEachElement(func(int, Element) error { return boom }); !errors.Is(err, boom) {
		t.Fatalf("iteration error=%v", err)
	}
	d.Pages[0].Number = 3
	if err := d.Validate(); err == nil {
		t.Fatal("bad number accepted")
	}
}
