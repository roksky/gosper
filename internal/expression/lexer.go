package expression

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxSourceBytes = 65536

type tokenKind uint8

const (
	tokEOF tokenKind = iota
	tokLiteral
	tokName
	tokReference
	tokOperator
	tokLParen
	tokRParen
)

type token struct {
	kind   tokenKind
	text   string
	value  Value
	offset int
}

func lex(src string) ([]token, error) {
	if len(src) > maxSourceBytes {
		return nil, fmt.Errorf("expression exceeds %d bytes", maxSourceBytes)
	}
	out := make([]token, 0, 16)
	for i := 0; i < len(src); {
		start := i
		r, size := utf8.DecodeRuneInString(src[i:])
		if r == utf8.RuneError && size == 1 {
			return nil, sourceError(src, i, "invalid UTF-8")
		}
		if unicode.IsSpace(r) {
			i += size
			continue
		}
		switch {
		case r == '$':
			if i+2 >= len(src) || (src[i+1] != 'F' && src[i+1] != 'P' && src[i+1] != 'V') || src[i+2] != '{' {
				return nil, sourceError(src, i, "invalid reference")
			}
			i += 3
			for i < len(src) && src[i] != '}' {
				i++
			}
			if i == len(src) {
				return nil, sourceError(src, start, "unterminated reference")
			}
			i++
			if _, err := ParseReference(src[start:i]); err != nil {
				return nil, sourceError(src, start, err.Error())
			}
			out = append(out, token{kind: tokReference, text: src[start:i], offset: start})
		case r == '"' || r == '\'':
			quote := byte(r)
			i++
			for i < len(src) {
				if src[i] == '\\' {
					i += 2
					continue
				}
				if src[i] == quote {
					break
				}
				_, n := utf8.DecodeRuneInString(src[i:])
				i += n
			}
			if i >= len(src) {
				return nil, sourceError(src, start, "unterminated string")
			}
			i++
			value, err := decodeString(src[start:i], quote)
			if err != nil {
				return nil, sourceError(src, start, "invalid string escape")
			}
			out = append(out, token{kind: tokLiteral, value: String(value), offset: start})
		case r >= '0' && r <= '9':
			i += size
			for i < len(src) && src[i] >= '0' && src[i] <= '9' {
				i++
			}
			decimal := false
			if i < len(src) && src[i] == '.' {
				decimal = true
				i++
				if i == len(src) || src[i] < '0' || src[i] > '9' {
					return nil, sourceError(src, i, "expected digit after decimal point")
				}
				for i < len(src) && src[i] >= '0' && src[i] <= '9' {
					i++
				}
			}
			raw := src[start:i]
			var v Value
			if decimal {
				v = Decimal(raw)
			} else {
				n, err := strconv.ParseInt(raw, 10, 64)
				if err != nil {
					return nil, sourceError(src, start, "integer literal overflows int64")
				}
				v = Int(n)
			}
			out = append(out, token{kind: tokLiteral, value: v, offset: start})
		case unicode.IsLetter(r) || r == '_':
			i += size
			for i < len(src) {
				rr, n := utf8.DecodeRuneInString(src[i:])
				if !unicode.IsLetter(rr) && !unicode.IsDigit(rr) && rr != '_' {
					break
				}
				i += n
			}
			word := src[start:i]
			switch word {
			case "true":
				out = append(out, token{kind: tokLiteral, value: Bool(true), offset: start})
			case "false":
				out = append(out, token{kind: tokLiteral, value: Bool(false), offset: start})
			case "null":
				out = append(out, token{kind: tokLiteral, value: NullValue, offset: start})
			case "and":
				out = append(out, token{kind: tokOperator, text: "&&", offset: start})
			case "or":
				out = append(out, token{kind: tokOperator, text: "||", offset: start})
			case "not":
				out = append(out, token{kind: tokOperator, text: "!", offset: start})
			default:
				out = append(out, token{kind: tokName, text: word, offset: start})
			}
		case r == '(':
			i++
			out = append(out, token{kind: tokLParen, offset: start})
		case r == ')':
			i++
			out = append(out, token{kind: tokRParen, offset: start})
		default:
			i += size
			op := src[start:i]
			if i < len(src) && (op == "=" || op == "!" || op == "<" || op == ">" || op == "&" || op == "|") && (src[i] == '=' || (src[i] == '&' && op == "&") || (src[i] == '|' && op == "|")) {
				op += src[i : i+1]
				i++
			}
			switch op {
			case "+", "-", "*", "/", "%", "==", "!=", "<", "<=", ">", ">=", "&&", "||", "!":
				out = append(out, token{kind: tokOperator, text: op, offset: start})
			default:
				return nil, sourceError(src, start, "unsupported character or operator")
			}
		}
		if len(out) > 8192 {
			return nil, sourceError(src, start, "too many tokens")
		}
	}
	out = append(out, token{kind: tokEOF, offset: len(src)})
	return out, nil
}

func decodeString(raw string, quote byte) (string, error) {
	if quote == '"' {
		return strconv.Unquote(raw)
	}
	content := raw[1 : len(raw)-1]
	var b strings.Builder
	for content != "" {
		r, _, tail, err := strconv.UnquoteChar(content, '\'')
		if err != nil {
			return "", err
		}
		b.WriteRune(r)
		content = tail
	}
	return b.String(), nil
}
