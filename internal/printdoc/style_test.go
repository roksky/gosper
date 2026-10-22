package printdoc

import "testing"

func TestDefaultAndExplicitStyles(t *testing.T) {
	d := DefaultStyle()
	if d.Font != "SansSerif" || d.Size != 200 || d.Color != 0 || d.Align != AlignLeft || d.Padding != (Insets{}) {
		t.Fatalf("default=%+v", d)
	}
	if err := d.Validate(); err != nil {
		t.Fatal(err)
	}
	e := NewTextElement(Bounds{Width: 100, Height: 30}, "x")
	if e.EffectiveStyle() != d {
		t.Fatalf("element default=%+v", e.EffectiveStyle())
	}
	explicit := Style{Font: "Serif", Size: 240, Bold: true, Italic: true, Color: 0x336699, Align: AlignRight, Padding: Insets{Top: 1, Right: 2, Bottom: 3, Left: 4}}
	e.Style = explicit
	if err := e.EffectiveStyle().Validate(); err != nil {
		t.Fatal(err)
	}
	if e.EffectiveStyle() != explicit {
		t.Fatalf("round trip=%+v", e.EffectiveStyle())
	}
}
func TestInvalidStyles(t *testing.T) {
	for _, s := range []Style{{Font: "Serif", Size: -1}, {Font: "Serif", Size: FromPoints(10001)}, {Font: "unknown", Size: 200}, {Font: "Serif", Size: 200, Color: 0x1000000}, {Font: "Serif", Size: 200, Align: Align(9)}, {Font: "Serif", Size: 200, Padding: Insets{Left: -1}}} {
		if err := s.Validate(); err == nil {
			t.Errorf("accepted %+v", s)
		}
	}
}
func TestDocumentStyleValidation(t *testing.T) {
	d := &Document{}
	e := NewTextElement(Bounds{Width: 100, Height: 20}, "hello")
	e.Style = Style{Font: "Serif", Size: 200, Color: 0x123456}
	if err := d.AppendPage(Page{Width: 1000, Height: 1000, Elements: []Element{e}}); err != nil {
		t.Fatal(err)
	}
	if d.Pages[0].Elements[0].Style != e.Style {
		t.Fatal("style changed")
	}
	e.Style.Color = 0xFFFFFFF
	if err := d.AppendPage(Page{Width: 1000, Height: 1000, Elements: []Element{e}}); err == nil {
		t.Fatal("invalid style accepted")
	}
}
