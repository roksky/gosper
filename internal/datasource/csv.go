package datasource

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/roksky/gosper/internal/diag"
)

// CSVOption configures a CSV source.
type CSVOption func(*csvSource)

// WithCSVDelimiter sets the field separator (comma by default).
func WithCSVDelimiter(delimiter rune) CSVOption {
	return func(s *csvSource) { s.delimiter = delimiter }
}

// FromCSV snapshots CSV input with a header row. Values remain strings;
// callers can convert them according to the report's declared field types.
func FromCSV(r io.Reader, opts ...CSVOption) DataSource {
	s := &csvSource{delimiter: ','}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	if r == nil {
		s.readErr = fmt.Errorf("nil reader")
	} else {
		s.data, s.readErr = io.ReadAll(r)
	}
	return s
}

type csvSource struct {
	data      []byte
	readErr   error
	delimiter rune
}

func (s *csvSource) Open(ctx context.Context) (Cursor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.readErr != nil {
		return nil, csvError(s.readErr)
	}
	data := bytes.TrimPrefix(s.data, []byte{0xef, 0xbb, 0xbf})
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = s.delimiter
	header, err := r.Read()
	if err != nil {
		return nil, csvError(err)
	}
	seen := make(map[string]bool, len(header))
	for _, name := range header {
		if strings.TrimSpace(name) == "" || seen[name] {
			return nil, csvError(fmt.Errorf("empty or duplicate header %q", name))
		}
		seen[name] = true
	}
	var rows []map[string]any
	for i := 2; ; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		fields, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, csvError(fmt.Errorf("row %d: %w", i, err))
		}
		row := make(map[string]any, len(header))
		for j, name := range header {
			row[name] = fields[j]
		}
		rows = append(rows, row)
	}
	return &csvCursor{rows: rows}, nil
}

func csvError(err error) error {
	d := diag.New(diag.CodeParse, diag.SeverityError, "CSV data source: "+err.Error(), nil)
	d.Cause = err
	return d
}

type csvCursor struct {
	rows   []map[string]any
	index  int
	closed bool
}

func (c *csvCursor) Next(ctx context.Context) (Record, bool, error) {
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

func (c *csvCursor) Close() error { c.closed = true; return nil }
