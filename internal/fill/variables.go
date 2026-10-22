package fill

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/roksky/gosper/internal/expression"
)

type Variable struct {
	Name, Calculation, ResetType, ResetGroup string
	Expr, Initial                            *expression.Program
}
type Group struct {
	Name string
	Expr *expression.Program
}
type variableState struct {
	values, groups map[string]expression.Value
	counts         map[string]int64
	started        bool
}

func newVariableState(t *Template) *variableState {
	return &variableState{values: map[string]expression.Value{}, groups: map[string]expression.Value{}, counts: map[string]int64{}}
}
func (s *variableState) apply(ctx context.Context, t *Template, record *FilledRecord) error {
	record.Groups = map[string]expression.Value{}
	changed := map[string]bool{}
	outerChanged := false
	for _, group := range t.Groups {
		value, err := group.Expr.Eval(ctx, expression.Scope{*record})
		if err != nil {
			return fmt.Errorf("group %s: %w", group.Name, err)
		}
		reset := !s.started || outerChanged || !reflect.DeepEqual(value, s.groups[group.Name])
		changed[group.Name] = reset
		outerChanged = outerChanged || reset
		if reset {
			s.counts[group.Name] = 0
		}
		s.counts[group.Name]++
		record.Variables[group.Name+"_COUNT"] = expression.Int(s.counts[group.Name])
		s.groups[group.Name] = value
		record.Groups[group.Name] = value
	}
	for _, v := range t.Variables {
		reset := !s.started || (v.ResetType == "Group" && changed[v.ResetGroup])
		if reset {
			initial := expression.Value(expression.NullValue)
			if v.Calculation == "Sum" || v.Calculation == "Count" {
				initial = expression.Int(0)
			}
			if v.Initial != nil {
				value, err := v.Initial.Eval(ctx, expression.Scope{*record})
				if err != nil {
					return err
				}
				initial = value
			}
			s.values[v.Name] = initial
		}
		for name, value := range s.values {
			record.Variables[name] = value
		}
		value := expression.Value(expression.NullValue)
		if v.Expr != nil {
			var err error
			value, err = v.Expr.Eval(ctx, expression.Scope{*record})
			if err != nil {
				return fmt.Errorf("variable %s: %w", v.Name, err)
			}
		}
		switch v.Calculation {
		case "Nothing", "":
			s.values[v.Name] = value
		case "Count":
			if value != nil && value != expression.NullValue {
				count, err := expression.AddNumbers(s.values[v.Name], expression.Int(1))
				if err != nil {
					return fmt.Errorf("variable %s: %w", v.Name, err)
				}
				s.values[v.Name] = count
			}
		case "Sum":
			if value != nil && value != expression.NullValue {
				sum, err := expression.AddNumbers(s.values[v.Name], value)
				if err != nil {
					return fmt.Errorf("variable %s: %w", v.Name, err)
				}
				s.values[v.Name] = sum
			}
		default:
			return fmt.Errorf("unsupported variable calculation %s", v.Calculation)
		}
		record.Variables[v.Name] = s.values[v.Name]
	}
	s.started = true
	return nil
}

func typedValue(raw any, class string) (expression.Value, error) {
	if raw == nil {
		return expression.NullValue, nil
	}
	if b, ok := raw.([]byte); ok {
		raw = string(b)
	}
	switch class {
	case "java.lang.String", "java.util.UUID":
		return expression.String(fmt.Sprint(raw)), nil
	case "java.math.BigDecimal", "java.lang.Double", "java.lang.Float":
		text := fmt.Sprint(raw)
		if _, err := strconv.ParseFloat(text, 64); err != nil {
			return nil, err
		}
		return expression.Decimal(text), nil
	case "java.lang.Integer", "java.lang.Long", "java.lang.Short":
		text := fmt.Sprint(raw)
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return nil, err
		}
		return expression.Int(v), nil
	case "java.sql.Timestamp", "java.sql.Date", "java.util.Date":
		if value, ok := raw.(time.Time); ok {
			return expression.Time(value), nil
		}
		text := strings.TrimSpace(fmt.Sprint(raw))
		for _, format := range []string{time.RFC3339Nano, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05.999999999", "2006-01-02"} {
			if value, err := time.Parse(format, text); err == nil {
				return expression.Time(value), nil
			}
		}
		return nil, fmt.Errorf("invalid date %q", text)
	}
	return ValueOf(raw)
}
