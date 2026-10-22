package datasource

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FromSlice creates a source from rows of structs, string-keyed maps, or
// positional slices/arrays. The top-level slice is copied; row values are not.
// Callers must not mutate rows while a cursor reads them.
func FromSlice[T any](rows []T) DataSource {
	copyRows := append([]T(nil), rows...)
	return &sliceSource[T]{rows: copyRows}
}

type sliceSource[T any] struct{ rows []T }

func (s *sliceSource[T]) Open(ctx context.Context) (Cursor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &sliceCursor[T]{rows: s.rows}, nil
}

type sliceCursor[T any] struct {
	rows   []T
	index  int
	closed bool
}

func (c *sliceCursor[T]) Next(ctx context.Context) (Record, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	if c.closed || c.index >= len(c.rows) {
		return nil, false, nil
	}
	row, err := adaptRow(c.rows[c.index])
	if err != nil {
		return nil, false, fmt.Errorf("row %d: %w", c.index, err)
	}
	c.index++
	return row, true, nil
}
func (c *sliceCursor[T]) Close() error { c.closed = true; return nil }

type sliceRecord map[string]any

func (r sliceRecord) Value(name string) (any, bool) { v, ok := r[name]; return v, ok }
func adaptRow(row any) (Record, error) {
	v := reflect.ValueOf(row)
	for v.IsValid() && (v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer) {
		if v.IsNil() {
			return sliceRecord{}, nil
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return sliceRecord{}, nil
	}
	out := sliceRecord{}
	switch v.Kind() {
	case reflect.Struct:
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			f := typ.Field(i)
			if !f.IsExported() {
				continue
			}
			name := f.Tag.Get("jrxml")
			if name == "" {
				name = f.Tag.Get("json")
			}
			if name != "" {
				name = strings.Split(name, ",")[0]
			}
			if name == "-" {
				continue
			}
			if name == "" {
				name = f.Name
			}
			out[name] = v.Field(i).Interface()
		}
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map key must be string")
		}
		iter := v.MapRange()
		for iter.Next() {
			out[iter.Key().String()] = iter.Value().Interface()
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			out[strconv.Itoa(i)] = v.Index(i).Interface()
		}
	default:
		return nil, fmt.Errorf("unsupported row type %s", v.Type())
	}
	return out, nil
}
