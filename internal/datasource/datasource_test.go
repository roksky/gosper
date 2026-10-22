package datasource

import (
	"context"
	"errors"
	"testing"
)

// Compile-time checks that the test fakes satisfy the contract.
var (
	_ DataSource = (*fakeSource)(nil)
	_ Cursor     = (*fakeCursor)(nil)
	_ Record     = (fakeRecord)(nil)
)

type fakeRecord map[string]any

func (r fakeRecord) Value(field string) (any, bool) {
	v, ok := r[field]
	return v, ok
}

type fakeCursor struct {
	recs     []fakeRecord
	idx      int
	closed   bool
	closeN   int
	closeErr error
	nextErr  error
}

func (c *fakeCursor) Next(ctx context.Context) (Record, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	if c.closed {
		return nil, false, nil
	}
	if c.nextErr != nil {
		return nil, false, c.nextErr
	}
	if c.idx >= len(c.recs) {
		return nil, false, nil
	}
	r := c.recs[c.idx]
	c.idx++
	return r, true, nil
}

func (c *fakeCursor) Close() error {
	c.closeN++
	if c.closed {
		return nil
	}
	c.closed = true
	return c.closeErr
}

type fakeSource struct {
	cursor  *fakeCursor
	openErr error
}

func (s *fakeSource) Open(ctx context.Context) (Cursor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.openErr != nil {
		return nil, s.openErr
	}
	return s.cursor, nil
}

func TestIterateAll(t *testing.T) {
	src := &fakeSource{cursor: &fakeCursor{recs: []fakeRecord{
		{"name": "a", "n": 1},
		{"name": "b", "n": 2},
		{"name": "c", "n": 3},
	}}}

	cur, err := src.Open(context.Background())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer cur.Close()

	var names []string
	for {
		rec, ok, err := cur.Next(context.Background())
		if err != nil {
			t.Fatalf("Next: %v", err)
		}
		if !ok {
			break
		}
		name, found := rec.Value("name")
		if !found {
			t.Fatalf("field %q missing", "name")
		}
		names = append(names, name.(string))
	}

	if want := []string{"a", "b", "c"}; len(names) != len(want) {
		t.Fatalf("names = %v, want %v", names, want)
	} else {
		for i := range want {
			if names[i] != want[i] {
				t.Fatalf("names = %v, want %v", names, want)
			}
		}
	}
}

func TestRecordMissingField(t *testing.T) {
	rec := fakeRecord{"present": 1}
	if v, ok := rec.Value("missing"); ok || v != nil {
		t.Fatalf("Value(missing) = (%v, %v), want (nil, false)", v, ok)
	}
	if v, ok := rec.Value("present"); !ok || v != 1 {
		t.Fatalf("Value(present) = (%v, %v), want (1, true)", v, ok)
	}
}

func TestNextAfterClose(t *testing.T) {
	cur := &fakeCursor{recs: []fakeRecord{{"a": 1}}}
	if err := cur.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	rec, ok, err := cur.Next(context.Background())
	if rec != nil || ok || err != nil {
		t.Fatalf("Next after Close = (%v, %v, %v), want (nil, false, nil)", rec, ok, err)
	}
}

func TestDoubleClose(t *testing.T) {
	cur := &fakeCursor{}
	if err := cur.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := cur.Close(); err != nil {
		t.Fatalf("second Close: %v, want nil", err)
	}
	if cur.closeN != 2 {
		t.Fatalf("Close called %d times, want 2", cur.closeN)
	}
}

func TestCloseAfterNextError(t *testing.T) {
	wantErr := errors.New("boom")
	cur := &fakeCursor{nextErr: wantErr}
	if _, _, err := cur.Next(context.Background()); !errors.Is(err, wantErr) {
		t.Fatalf("Next error = %v, want %v", err, wantErr)
	}
	if err := cur.Close(); err != nil {
		t.Fatalf("Close after Next error = %v, want nil", err)
	}
}

func TestOpenCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	src := &fakeSource{cursor: &fakeCursor{}}
	if _, err := src.Open(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Open with canceled ctx = %v, want context.Canceled", err)
	}
}

func TestNextCancellation(t *testing.T) {
	cur := &fakeCursor{recs: []fakeRecord{{"a": 1}}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, _, err := cur.Next(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Next with canceled ctx = %v, want context.Canceled", err)
	}
}

func TestOpenError(t *testing.T) {
	wantErr := errors.New("open failed")
	src := &fakeSource{openErr: wantErr}
	if _, err := src.Open(context.Background()); !errors.Is(err, wantErr) {
		t.Fatalf("Open error = %v, want %v", err, wantErr)
	}
}
