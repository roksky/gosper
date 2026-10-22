package report

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/limits"
)

// Query resolves parameter defaults and binds the embedded SQL query. Values
// are always returned separately; raw $P!{} and $X{} clauses are unsupported.
// placeholder receives a one-based argument index (nil produces ?).
func (c *CompiledReport) Query(ctx context.Context, params Params, placeholder func(int) string) (string, []any, error) {
	if c == nil {
		return "", nil, fmt.Errorf("report: nil compiled report")
	}
	overrides := make(map[string]expression.Value, len(params))
	for name, raw := range params {
		v, err := fill.ValueOf(raw)
		if err != nil {
			return "", nil, fmt.Errorf("parameter %q: %w", name, err)
		}
		overrides[name] = v
	}
	resolved, err := fill.RunWithParams(ctx, &c.fields, FromSlice([]map[string]any{}), overrides, limits.Budget{})
	if err != nil {
		return "", nil, err
	}
	values := make(Params, len(resolved.Parameters))
	for name, v := range resolved.Parameters {
		switch x := v.(type) {
		case expression.Null:
			values[name] = nil
		case expression.String:
			values[name] = string(x)
		case expression.Int:
			values[name] = int64(x)
		case expression.Bool:
			values[name] = bool(x)
		case expression.Decimal:
			values[name] = string(x)
		case expression.Time:
			values[name] = time.Time(x)
		}
	}
	return BindQuery(c.query, values, placeholder)
}

// BindQuery replaces Jasper $P{name} references with driver placeholders.
// SQL string literals, quoted identifiers and comments are preserved verbatim.
func BindQuery(query string, params Params, placeholder func(int) string) (string, []any, error) {
	if placeholder == nil {
		placeholder = func(int) string { return "?" }
	}
	var out strings.Builder
	var args []any
	for i := 0; i < len(query); {
		start := i
		switch {
		case strings.HasPrefix(query[i:], "--"):
			end := strings.IndexByte(query[i:], '\n')
			if end < 0 {
				end = len(query) - i
			}
			i += end
		case strings.HasPrefix(query[i:], "/*"):
			depth := 1
			i += 2
			for i < len(query) && depth > 0 {
				if strings.HasPrefix(query[i:], "/*") {
					depth++
					i += 2
				} else if strings.HasPrefix(query[i:], "*/") {
					depth--
					i += 2
				} else {
					i++
				}
			}
			if depth > 0 {
				return "", nil, fmt.Errorf("report: unterminated SQL comment")
			}
		case query[i] == '\'' || query[i] == '"':
			quote := query[i]
			i++
			closed := false
			for i < len(query) {
				if query[i] == quote {
					i++
					if i < len(query) && query[i] == quote {
						i++
						continue
					}
					closed = true
					break
				}
				i++
			}
			if !closed {
				return "", nil, fmt.Errorf("report: unterminated SQL quote")
			}
		case query[i] == '$' && dollarQuote.FindString(query[i:]) != "":
			delimiter := dollarQuote.FindString(query[i:])
			end := strings.Index(query[i+len(delimiter):], delimiter)
			if end < 0 {
				return "", nil, fmt.Errorf("report: unterminated SQL dollar quote")
			}
			i += len(delimiter)*2 + end
		case strings.HasPrefix(query[i:], "$P!{") || strings.HasPrefix(query[i:], "$X{"):
			return "", nil, fmt.Errorf("report: raw parameters and query clauses: %w", ErrNotImplemented)
		case strings.HasPrefix(query[i:], "$P{"):
			end := strings.IndexByte(query[i+3:], '}')
			if end < 0 {
				return "", nil, fmt.Errorf("report: unterminated query parameter")
			}
			name := query[i+3 : i+3+end]
			value, ok := params[name]
			if !ok {
				return "", nil, fmt.Errorf("report: missing query parameter %q", name)
			}
			args = append(args, value)
			out.WriteString(placeholder(len(args)))
			i += end + 4
			continue
		default:
			i++
		}
		out.WriteString(query[start:i])
	}
	return out.String(), args, nil
}

var dollarQuote = regexp.MustCompile(`^\$(?:[a-zA-Z_][a-zA-Z0-9_]*)?\$`)
