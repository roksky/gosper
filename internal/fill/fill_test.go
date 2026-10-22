package fill

import (
	"context"
	"errors"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/limits"
	"sync"
	"testing"
)

type testRecord map[string]any

func (r testRecord) Value(n string) (any, bool) { v, ok := r[n]; return v, ok }

type testCursor struct {
	rows     []testRecord
	idx      int
	closed   bool
	nextErr  error
	closeErr error
	cancel   context.CancelFunc
}

func (c *testCursor) Next(ctx context.Context) (datasource.Record, bool, error) {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	if c.nextErr != nil {
		return nil, false, c.nextErr
	}
	if c.idx >= len(c.rows) {
		return nil, false, nil
	}
	r := c.rows[c.idx]
	c.idx++
	return r, true, nil
}
func (c *testCursor) Close() error { c.closed = true; return c.closeErr }

type testSource struct{ cur *testCursor }

func (s testSource) Open(ctx context.Context) (datasource.Cursor, error) { return s.cur, nil }
func TestRunSnapshotsAndMissing(t *testing.T) {
	src := datasource.FromSlice([]map[string]any{{"n": 1, "name": "a"}, {"n": 2}, {"n": 3, "name": nil}})
	got, err := Run(context.Background(), &Template{Fields: []string{"n", "name"}}, src, limits.Budget{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 3 {
		t.Fatalf("records=%d", len(got.Records))
	}
	if v, ok := got.Records[0].Field("n"); !ok || v != expression.Int(1) {
		t.Errorf("first n=%v,%v", v, ok)
	}
	if _, ok := got.Records[1].Field("name"); ok {
		t.Error("missing field present")
	}
	if v, ok := got.Records[2].Field("name"); !ok || v != expression.NullValue {
		t.Errorf("null field=%v,%v", v, ok)
	}
}
func TestRunClosesOnFailure(t *testing.T) {
	boom := errors.New("next failed")
	for _, tc := range []struct {
		name string
		cur  *testCursor
		want error
	}{{"next", &testCursor{nextErr: boom}, boom}, {"conversion", &testCursor{rows: []testRecord{{"n": struct{}{}}}}, nil}, {"close", &testCursor{closeErr: boom}, boom}} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Run(context.Background(), &Template{Fields: []string{"n"}}, testSource{tc.cur}, limits.Budget{})
			if err == nil {
				t.Fatal("expected error")
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("error=%v", err)
			}
			if !tc.cur.closed {
				t.Fatal("cursor left open")
			}
		})
	}
}
func TestRunCancellationAndLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cur := &testCursor{rows: []testRecord{{"n": 1}}, cancel: cancel}
	_, err := Run(ctx, &Template{Fields: []string{"n"}}, testSource{cur}, limits.Budget{})
	if !errors.Is(err, context.Canceled) || !cur.closed {
		t.Fatalf("cancel=%v, closed=%v", err, cur.closed)
	}
	cur = &testCursor{rows: []testRecord{{"n": 1}, {"n": 2}}}
	_, err = Run(context.Background(), &Template{Fields: []string{"n"}}, testSource{cur}, limits.Budget{MaxRecords: 1})
	if err == nil || !cur.closed {
		t.Fatalf("limit=%v, closed=%v", err, cur.closed)
	}
}
func TestRunConcurrentIsolation(t *testing.T) {
	src := datasource.FromSlice([]map[string]int{{"n": 1}, {"n": 2}})
	tmpl := &Template{Fields: []string{"n"}}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := Run(context.Background(), tmpl, src, limits.Budget{})
			if err != nil || len(got.Records) != 2 {
				t.Errorf("run=%v,%v", got, err)
			}
		}()
	}
	wg.Wait()
}
