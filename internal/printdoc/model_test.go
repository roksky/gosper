package printdoc

import (
	"math"
	"testing"
)

func TestUnitConversions(t *testing.T) {
	cases := []struct {
		name string
		unit Unit
		want float64
		got  func(Unit) float64
	}{
		{"1 inch in points", FromInches(1), 72, Unit.Points},
		{"1 inch in inches", FromInches(1), 1, Unit.Inches},
		{"1 inch in mm", FromInches(1), 25.4, Unit.MM},
		{"72 points in inches", FromPoints(72), 1, Unit.Inches},
		{"25.4mm in inches", FromMM(25.4), 1, Unit.Inches},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.got(tc.unit); math.Abs(got-tc.want) > 1e-9 {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUnitExactValues(t *testing.T) {
	if got := FromInches(1); got != 1440 {
		t.Fatalf("FromInches(1) = %d, want 1440", got)
	}
	if got := FromPoints(1); got != 20 {
		t.Fatalf("FromPoints(1) = %d, want 20", got)
	}
	// Rounding is half away from zero.
	if got := FromPoints(0.1); got != 2 {
		t.Fatalf("FromPoints(0.1) = %d, want 2", got)
	}
}

func TestNoCumulativeError(t *testing.T) {
	// Ten times 0.1pt must equal exactly 1pt worth of twips.
	var sum Unit
	for i := 0; i < 10; i++ {
		sum += FromPoints(0.1)
	}
	if sum != FromPoints(1) {
		t.Fatalf("sum of ten 0.1pt = %d, want %d", sum, FromPoints(1))
	}
}

func TestPageValidate(t *testing.T) {
	ok := Page{Number: 1, Width: FromInches(8.5), Height: FromInches(11)}
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid page rejected: %v", err)
	}
	for _, bad := range []Page{
		{Number: 1, Width: 0, Height: FromInches(11)},
		{Number: 2, Width: FromInches(8.5), Height: -1},
	} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("page %+v should be invalid", bad)
		}
	}
}

func TestDocumentModel(t *testing.T) {
	doc := Document{
		Pages: []Page{
			{
				Number: 1,
				Width:  FromInches(8.5),
				Height: FromInches(11),
				Elements: []Element{{
					Kind:   ElementText,
					Bounds: Bounds{X: FromInches(1), Y: FromInches(1), Width: FromInches(2), Height: FromPoints(12)},
					Text:   "hello",
				}},
			},
		},
		Metadata: Metadata{PageCount: 1},
	}

	if doc.Metadata.PageCount != len(doc.Pages) {
		t.Fatalf("PageCount %d != len(Pages) %d", doc.Metadata.PageCount, len(doc.Pages))
	}
	if got := doc.Pages[0].Elements[0].Bounds.X; got != 1440 {
		t.Fatalf("element X = %d, want 1440", got)
	}
	if got := doc.Pages[0].Elements[0].Kind.String(); got != "text" {
		t.Fatalf("kind = %q, want text", got)
	}
}

func TestLargePageNoOverflow(t *testing.T) {
	// 200 inches is 288000 twips, well within int64.
	p := Page{Number: 1, Width: FromInches(200), Height: FromInches(200)}
	if err := p.Validate(); err != nil {
		t.Fatalf("large page rejected: %v", err)
	}
	if p.Width != 288000 {
		t.Fatalf("width = %d, want 288000", p.Width)
	}
}
