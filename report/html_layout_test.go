package report

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/printdoc"
)

func TestHTMLPreservesPrintPositionsAndEmbeddedImages(t *testing.T) {
	var pngData bytes.Buffer
	logo := image.NewRGBA(image.Rect(0, 0, 2, 2))
	logo.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(&pngData, logo); err != nil {
		t.Fatal(err)
	}
	resource := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngData.Bytes())
	text := printdoc.NewTextElement(printdoc.Bounds{X: printdoc.FromPoints(40), Y: printdoc.FromPoints(30), Width: printdoc.FromPoints(10), Height: printdoc.FromPoints(10)}, "Wide <Alice> & Bob")
	shape := printdoc.NewShapeElement(printdoc.Bounds{Width: printdoc.FromPoints(100), Height: printdoc.FromPoints(1)}, printdoc.ShapeLine)
	shape.Shape.Stroke = true
	doc := &printdoc.Document{}
	if err := doc.AppendPage(printdoc.Page{Width: printdoc.FromPoints(595), Height: printdoc.FromPoints(842), Elements: []printdoc.Element{shape, text, printdoc.NewImageElement(printdoc.Bounds{Width: printdoc.FromPoints(50), Height: printdoc.FromPoints(50)}, resource)}}); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := Export(context.Background(), &FilledReport{document: doc}, &output, FormatHTML); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, want := range []string{"html,body{margin:0;padding:0", `left:40pt;top:30pt;width:10pt;height:10pt`, `<text x="0" y="8"`, `white-space:pre;`, `Wide &lt;Alice&gt; &amp; Bob`, `<line x1="0" y1="0" x2="100" y2="1"`, resource, "object-fit:contain"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in %s", want, html)
		}
	}
	if strings.Contains(html, "white-space:pre-wrap") {
		t.Fatal("HTML re-wraps engine-positioned text")
	}
}
