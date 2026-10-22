package datasource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/roksky/gosper/internal/diag"
)

// FromJSON snapshots a JSON array of objects. Each Open decodes an independent
// cursor. Numbers are float64, and nested objects and arrays remain Go maps and
// slices. A present null field returns (nil, true).
func FromJSON(r io.Reader) DataSource {
	if r == nil {
		return &jsonSource{readErr: fmt.Errorf("nil reader")}
	}
	data, err := io.ReadAll(r)
	return &jsonSource{data: data, readErr: err}
}

type jsonSource struct {
	data    []byte
	readErr error
}

func (s *jsonSource) Open(ctx context.Context) (Cursor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.readErr != nil {
		return nil, jsonError(s.readErr)
	}
	dec := json.NewDecoder(bytes.NewReader(s.data))
	var rows []map[string]any
	if err := dec.Decode(&rows); err != nil {
		return nil, jsonError(err)
	}
	if rows == nil || len(bytes.TrimSpace(s.data)) == 0 || bytes.TrimSpace(s.data)[0] != '[' {
		return nil, jsonError(fmt.Errorf("top-level value must be an array of objects"))
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("trailing JSON value")
		}
		return nil, jsonError(err)
	}
	for i, row := range rows {
		if row == nil {
			return nil, jsonError(fmt.Errorf("row %d must be an object", i+1))
		}
	}
	return &jsonCursor{rows: rows}, nil
}

func jsonError(err error) error {
	d := diag.New(diag.CodeParse, diag.SeverityError, "JSON data source: "+err.Error(), nil)
	d.Cause = err
	return d
}

type jsonCursor struct {
	rows   []map[string]any
	index  int
	closed bool
}

func (c *jsonCursor) Next(ctx context.Context) (Record, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	if c.closed || c.index == len(c.rows) {
		return nil, false, nil
	}
	row := sliceRecord(c.rows[c.index])
	c.index++
	return row, true, nil
}

func (c *jsonCursor) Close() error { c.closed = true; return nil }
