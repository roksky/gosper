package jrxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/roksky/gosper/internal/diag"
)

// handlerFunc parses one child element of <jasperReport>.
//
// Implementations consume tokens from p.dec starting immediately after se
// (which has already been consumed) up to and including its matching
// EndElement. They append to p.tmpl and/or p.diags. A returned error is
// treated as a fatal XML/IO failure.
type handlerFunc func(p *parser, se xml.StartElement) error

// elementParser parses one drawable element found inside a band.
type elementParser func(p *parser, se xml.StartElement) (Element, error)

// Registration seams. Handler files call registerRootChild / registerChildElement
// / registerRootStart from init, so new constructs are added in new files without
// editing the dispatcher. Duplicate registrations are a programming error.
var (
	rootChildHandlers   = map[string]handlerFunc{}
	childElementParsers = map[string]elementParser{}
	rootStartHandlers   []func(p *parser, se xml.StartElement)
)

func registerRootChild(name string, h handlerFunc) {
	if _, dup := rootChildHandlers[name]; dup {
		panic("jrxml: duplicate root handler for <" + name + ">")
	}
	rootChildHandlers[name] = h
}

func registerChildElement(name string, fn elementParser) {
	if _, dup := childElementParsers[name]; dup {
		panic("jrxml: duplicate element parser for <" + name + ">")
	}
	childElementParsers[name] = fn
}

func registerRootStart(f func(p *parser, se xml.StartElement)) {
	rootStartHandlers = append(rootStartHandlers, f)
}

// parser carries the decoder, source bytes (for spans) and accumulated results.
type parser struct {
	dec   *xml.Decoder
	data  []byte
	tmpl  *Template
	diags diag.List

	// seenSections tracks band sections already parsed, for duplicate detection.
	seenSections map[BandKind]bool
}

// Parse reads a JRXML template from r.
//
// It returns the parsed Template and a list of non-fatal diagnostics. On a fatal
// error (unreadable input, malformed XML, or no root element) the returned
// Template is nil, the returned list contains a diagnostic with code
// diag.CodeParse, and the returned error wraps that diagnostic so callers can
// recover it with errors.As or diag.As. Parse never panics.
func Parse(r io.Reader) (*Template, diag.List, error) {
	if r == nil {
		return fail(diag.CodeInternal, "nil reader", nil)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return fail(diag.CodeInternal, "reading template: "+err.Error(), nil)
	}
	return parse(data)
}

func parse(data []byte) (*Template, diag.List, error) {
	// A leading UTF-8 BOM is not valid XML content; strip it so encoding/xml
	// sees the declaration or root element directly.
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	dec := xml.NewDecoder(bytes.NewReader(data))
	root, ok, err := firstStartElement(dec)
	if err != nil {
		return fail(diag.CodeParse, "malformed XML: "+err.Error(), syntaxSpan(err))
	}
	if !ok {
		return fail(diag.CodeParse, "no root element", &diag.Span{Line: 1, Col: 1})
	}

	p := &parser{
		dec:          dec,
		data:         data,
		tmpl:         &Template{Page: defaultPage()},
		seenSections: map[BandKind]bool{},
	}

	version, _ := attribute(root, "version")
	p.tmpl.Dialect = parseDialect(version)
	if !p.tmpl.Dialect.Known() {
		p.diags.Add(diag.New(
			diag.CodeUnsupported,
			diag.SeverityWarning,
			fmt.Sprintf("unsupported or missing JRXML version %q", version),
			spanForOffset(data, 0),
		))
	}

	if root.Name.Local != "jasperReport" {
		p.unsupported(root, root.Name.Local)
	}

	for _, f := range rootStartHandlers {
		f(p, root)
	}

	if err := p.parseRootChildren(); err != nil {
		return fail(diag.CodeParse, "malformed XML: "+err.Error(), syntaxSpan(err))
	}
	if err := p.drain(); err != nil {
		return fail(diag.CodeParse, "malformed XML: "+err.Error(), syntaxSpan(err))
	}
	return p.tmpl, p.diags, nil
}

// parseRootChildren consumes events up to and including the EndElement that
// closes the root.
func (p *parser) parseRootChildren() error {
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if err := p.dispatchRootChild(t); err != nil {
				return err
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (p *parser) dispatchRootChild(se xml.StartElement) error {
	if h, ok := rootChildHandlers[se.Name.Local]; ok {
		return h(p, se)
	}
	p.unsupported(se, se.Name.Local)
	return p.skipElement()
}

// skipElement consumes tokens until the EndElement matching the already-read
// StartElement. It must be called immediately after reading a StartElement.
func (p *parser) skipElement() error {
	depth := 0
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			if depth == 0 {
				return nil
			}
			depth--
		}
	}
}

// readText consumes the current element and returns its concatenated character
// data verbatim (no trimming). Nested elements are skipped.
func (p *parser) readText() (string, error) {
	var b strings.Builder
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return "", err
		}
		switch t := tok.(type) {
		case xml.CharData:
			b.Write(t)
		case xml.StartElement:
			if err := p.skipElement(); err != nil {
				return "", err
			}
		case xml.EndElement:
			return b.String(), nil
		}
	}
}

// drain reads tokens until EOF to validate well-formedness of any trailing
// content after the root element.
func (p *parser) drain() error {
	for {
		_, err := p.dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// firstStartElement returns the first start element token, skipping the XML
// declaration, comments, directives and whitespace. ok is false if the document
// contains no start element.
func firstStartElement(dec *xml.Decoder) (xml.StartElement, bool, error) {
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return xml.StartElement{}, false, nil
		}
		if err != nil {
			return xml.StartElement{}, false, err
		}
		if se, ok := tok.(xml.StartElement); ok {
			return se, true, nil
		}
	}
}

// warnValidation appends a non-fatal validation warning.
func (p *parser) warnValidation(msg string) {
	p.diags.Add(diag.New(diag.CodeValidation, diag.SeverityWarning, msg, p.span()))
}

// span returns a best-effort span at the decoder's current input offset. It is
// never nil, so diagnostics always carry a location.
func (p *parser) span() *diag.Span {
	return spanForOffset(p.data, int(p.dec.InputOffset()))
}

// syntaxSpan derives a best-effort span from an encoding/xml error. Line comes
// from *xml.SyntaxError; column and offset are unknown and left at zero.
func syntaxSpan(err error) *diag.Span {
	var se *xml.SyntaxError
	if errors.As(err, &se) {
		return &diag.Span{Line: se.Line}
	}
	return nil
}

// spanForOffset computes a 1-based line/column span for a byte offset into data.
func spanForOffset(data []byte, off int) *diag.Span {
	if off < 0 {
		off = 0
	}
	if off > len(data) {
		off = len(data)
	}
	line, col := 1, 1
	for i := 0; i < off; i++ {
		if data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return &diag.Span{Line: line, Col: col, Offset: off}
}

// fail builds a fatal result: a nil template, a single-diagnostic list, and an
// error wrapping that diagnostic.
func fail(code diag.Code, msg string, span *diag.Span) (*Template, diag.List, error) {
	var l diag.List
	l.Add(diag.New(code, diag.SeverityError, msg, span))
	return nil, l, l.Err()
}
