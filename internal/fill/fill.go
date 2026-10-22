package fill

import (
	"context"
	"errors"
	"fmt"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/limits"
)

// Template is the immutable field selection for a fill run. Callers must not
// mutate Fields while a run is active.
type Template struct {
	Fields     []string
	Parameters []ParameterDefault
}

// FilledRecord owns a snapshot of the requested fields for one source record.
type FilledRecord struct {
	Fields     map[string]expression.Value
	Parameters map[string]expression.Value
	Variables  map[string]expression.Value
}

// Field implements the field part of expression.Symbols.
func (r FilledRecord) Field(name string) (expression.Value, bool) {
	v, ok := r.Fields[name]
	return v, ok
}

// Parameter resolves a parameter snapshot for this record.
func (r FilledRecord) Parameter(name string) (expression.Value, bool) {
	v, ok := r.Parameters[name]
	return v, ok
}

// Variable resolves a built-in variable snapshot for this record.
func (r FilledRecord) Variable(name string) (expression.Value, bool) {
	v, ok := r.Variables[name]
	return v, ok
}

// Filled holds the record snapshots produced by one fill run.
type Filled struct {
	Records    []FilledRecord
	Parameters map[string]expression.Value
}

// Run opens the source, copies selected fields from each record, and closes
// the cursor on every exit path. Each invocation owns all mutable run state.
func Run(ctx context.Context, tmpl *Template, ds datasource.DataSource, budget limits.Budget) (result *Filled, err error) {
	return RunWithParams(ctx, tmpl, ds, nil, budget)
}

// RunWithParams runs with supplied parameter overrides and declared defaults.
func RunWithParams(ctx context.Context, tmpl *Template, ds datasource.DataSource, overrides map[string]expression.Value, budget limits.Budget) (result *Filled, err error) {
	if ctx == nil {
		return nil, fmt.Errorf("fill: nil context")
	}
	if tmpl == nil || ds == nil {
		return nil, fmt.Errorf("fill: nil template or data source")
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	params, err := resolveParameters(ctx, tmpl.Parameters, overrides)
	if err != nil {
		return nil, err
	}
	cur, err := ds.Open(ctx)
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, fmt.Errorf("fill: source returned nil cursor")
	}
	defer func() {
		err = errors.Join(err, cur.Close())
		if err != nil {
			result = nil
		}
	}()
	guard := budget.Guard(ctx)
	filled := &Filled{Parameters: cloneValues(params)}
	for {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		var rec datasource.Record
		var ok bool
		rec, ok, err = cur.Next(ctx)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		if rec == nil {
			return nil, fmt.Errorf("fill: source returned nil record")
		}
		if err = guard.Record(); err != nil {
			return nil, err
		}
		fields := make(map[string]expression.Value, len(tmpl.Fields))
		for _, name := range tmpl.Fields {
			raw, found := rec.Value(name)
			if !found {
				continue
			}
			var value expression.Value
			value, err = ValueOf(raw)
			if err != nil {
				return nil, fmt.Errorf("field %q: %w", name, err)
			}
			fields[name] = value
		}
		filled.Records = append(filled.Records, FilledRecord{
			Fields: fields, Parameters: cloneValues(params),
			Variables: map[string]expression.Value{"REPORT_COUNT": expression.Int(len(filled.Records) + 1)},
		})
	}
	return filled, nil
}

// ValueOf converts supported Go field and parameter values to expression values.
func ValueOf(raw any) (expression.Value, error) {
	if raw == nil {
		return expression.NullValue, nil
	}
	if value, ok := raw.(expression.Value); ok {
		return value, nil
	}
	switch v := raw.(type) {
	case string:
		return expression.String(v), nil
	case bool:
		return expression.Bool(v), nil
	case int:
		return expression.Int(v), nil
	case int8:
		return expression.Int(v), nil
	case int16:
		return expression.Int(v), nil
	case int32:
		return expression.Int(v), nil
	case int64:
		return expression.Int(v), nil
	case uint, uint8, uint16, uint32, uint64:
		var n uint64
		switch x := raw.(type) {
		case uint:
			n = uint64(x)
		case uint8:
			n = uint64(x)
		case uint16:
			n = uint64(x)
		case uint32:
			n = uint64(x)
		case uint64:
			n = x
		}
		if n > uint64(^uint64(0)>>1) {
			return nil, fmt.Errorf("unsigned integer overflows int64")
		}
		return expression.Int(n), nil
	default:
		return nil, fmt.Errorf("unsupported value type %T", raw)
	}
}
