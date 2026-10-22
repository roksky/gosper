package expression

import (
	"fmt"
	"github.com/roksky/gosper/internal/diag"
)

const maxParseDepth = 128

type node struct {
	op          string
	value       Value
	ref         Reference
	left, right *node
	offset      int
}

// Program is an immutable parsed expression, safe to evaluate concurrently.
type Program struct {
	src  string
	root *node
}

type parser struct {
	src    string
	tokens []token
	pos    int
	depth  int
}

// Compile parses the closed expression grammar. It does not execute code or
// read external state. Decimal literals are accepted but arithmetic on them
// awaits the decimal backend.
func Compile(src string) (*Program, error) {
	tokens, err := lex(src)
	if err != nil {
		return nil, err
	}
	p := &parser{src: src, tokens: tokens}
	root, err := p.parse(1)
	if err != nil {
		return nil, err
	}
	if p.current().kind != tokEOF {
		return nil, sourceError(src, p.current().offset, "unexpected token")
	}
	return &Program{src: src, root: root}, nil
}
func (p *parser) current() token { return p.tokens[p.pos] }
func (p *parser) parse(min int) (*node, error) {
	p.depth++
	defer func() { p.depth-- }()
	if p.depth > maxParseDepth {
		return nil, sourceError(p.src, p.current().offset, "expression nesting too deep")
	}
	t := p.current()
	p.pos++
	var left *node
	switch t.kind {
	case tokLiteral:
		left = &node{value: t.value, offset: t.offset}
	case tokName, tokReference:
		ref, err := ParseReference(t.text)
		if err != nil {
			return nil, sourceError(p.src, t.offset, err.Error())
		}
		left = &node{ref: ref, op: "ref", offset: t.offset}
	case tokLParen:
		var err error
		left, err = p.parse(1)
		if err != nil {
			return nil, err
		}
		if p.current().kind != tokRParen {
			return nil, sourceError(p.src, p.current().offset, "expected ')'")
		}
		p.pos++
	case tokOperator:
		if t.text != "!" && t.text != "-" && t.text != "+" {
			return nil, sourceError(p.src, t.offset, "expected expression")
		}
		right, err := p.parse(7)
		if err != nil {
			return nil, err
		}
		left = &node{op: "u" + t.text, right: right, offset: t.offset}
	default:
		return nil, sourceError(p.src, t.offset, "expected expression")
	}
	for {
		t = p.current()
		if t.kind != tokOperator {
			break
		}
		prec := precedence(t.text)
		if prec < min {
			break
		}
		p.pos++
		right, err := p.parse(prec + 1)
		if err != nil {
			return nil, err
		}
		left = &node{op: t.text, left: left, right: right, offset: t.offset}
	}
	return left, nil
}
func precedence(op string) int {
	switch op {
	case "||":
		return 1
	case "&&":
		return 2
	case "==", "!=":
		return 3
	case "<", "<=", ">", ">=":
		return 4
	case "+", "-":
		return 5
	case "*", "/", "%":
		return 6
	}
	return 0
}
func sourceError(src string, offset int, message string) error {
	line, col := 1, 1
	for i := 0; i < offset && i < len(src); i++ {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return diag.New(diag.CodeParse, diag.SeverityError, message, &diag.Span{Line: line, Col: col, Offset: offset})
}
func runtimeError(src string, offset int, message string) error {
	return fmt.Errorf("%s", sourceError(src, offset, message))
}
