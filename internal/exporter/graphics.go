package exporter

import (
	"bytes"
	"compress/zlib"
	"encoding/ascii85"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"strings"

	"github.com/roksky/gosper/internal/printdoc"
)

func writeGraphic(out *countedWriter, page printdoc.Page, el printdoc.Element) error {
	b := el.Bounds
	x, y, w, h := b.X.Points(), (page.Height - b.Y - b.Height).Points(), b.Width.Points(), b.Height.Points()
	if el.Kind == printdoc.ElementImage {
		const prefix = "data:image/png;base64,"
		if !strings.HasPrefix(el.Image.Resource, prefix) {
			return fmt.Errorf("pdf: unsupported image resource")
		}
		data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(el.Image.Resource, prefix))
		if err != nil {
			return err
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
		rect := img.Bounds()
		ratio := float64(rect.Dx()) / float64(rect.Dy())
		if w/h > ratio {
			newWidth := h * ratio
			x += (w - newWidth) / 2
			w = newWidth
		} else {
			newHeight := w / ratio
			y += (h - newHeight) / 2
			h = newHeight
		}
		if rect.Dx() > 4096 || rect.Dy() > 4096 {
			return fmt.Errorf("pdf: image exceeds 4096 pixels")
		}
		var raw bytes.Buffer
		compressor := zlib.NewWriter(&raw)
		for iy := rect.Min.Y; iy < rect.Max.Y; iy++ {
			for ix := rect.Min.X; ix < rect.Max.X; ix++ {
				r, g, b, a := img.At(ix, iy).RGBA()
				pixel := []byte{byte((r + 65535 - a) >> 8), byte((g + 65535 - a) >> 8), byte((b + 65535 - a) >> 8)}
				if _, err := compressor.Write(pixel); err != nil {
					return err
				}
			}
		}
		if err := compressor.Close(); err != nil {
			return err
		}
		encoded := make([]byte, ascii85.MaxEncodedLen(raw.Len()))
		n := ascii85.Encode(encoded, raw.Bytes())
		// Inline images avoid shared mutable resources and preserve page streaming.
		if err := out.write(fmt.Sprintf("q %.4f 0 0 %.4f %.4f %.4f cm\nBI /W %d /H %d /BPC 8 /CS /RGB /F [/ASCII85Decode /FlateDecode] ID\n", w, h, x, y, rect.Dx(), rect.Dy())); err != nil {
			return err
		}
		if err := out.write(string(encoded[:n]) + "~>\nEI\nQ\n"); err != nil {
			return err
		}
		return nil
	}
	if el.Kind != printdoc.ElementShape {
		return fmt.Errorf("pdf: unsupported graphic %s", el.Kind)
	}
	s := el.Shape
	var command strings.Builder
	command.WriteString("q\n")
	rgb := func(color uint32, op string) {
		fmt.Fprintf(&command, "%.4f %.4f %.4f %s\n", float64((color>>16)&255)/255, float64((color>>8)&255)/255, float64(color&255)/255, op)
	}
	if s.Fill {
		rgb(s.FillColor, "rg")
	}
	if s.Stroke {
		rgb(s.StrokeColor, "RG")
		fmt.Fprintf(&command, "%.4f w\n", s.StrokeWidth.Points())
	}
	switch s.Kind {
	case printdoc.ShapeLine:
		fmt.Fprintf(&command, "%.4f %.4f m %.4f %.4f l\n", x, y+h, x+w, y)
	case printdoc.ShapeRectangle:
		fmt.Fprintf(&command, "%.4f %.4f %.4f %.4f re\n", x, y, w, h)
	case printdoc.ShapeEllipse:
		k := 0.552284749831
		cx, cy := x+w/2, y+h/2
		rx, ry := w/2, h/2
		fmt.Fprintf(&command, "%f %f m %f %f %f %f %f %f c %f %f %f %f %f %f c %f %f %f %f %f %f c %f %f %f %f %f %f c\n", cx+rx, cy, cx+rx, cy+k*ry, cx+k*rx, cy+ry, cx, cy+ry, cx-k*rx, cy+ry, cx-rx, cy+k*ry, cx-rx, cy, cx-rx, cy-k*ry, cx-k*rx, cy-ry, cx, cy-ry, cx+k*rx, cy-ry, cx+rx, cy-k*ry, cx+rx, cy)
	}
	if s.Fill && s.Stroke {
		command.WriteString("B\n")
	} else if s.Fill {
		command.WriteString("f\n")
	} else {
		command.WriteString("S\n")
	}
	command.WriteString("Q\n")
	return out.write(command.String())
}
