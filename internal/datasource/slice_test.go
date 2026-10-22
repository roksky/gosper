package datasource

import (
	"context"
	"errors"
	"testing"
)

func TestSliceSourceRows(t *testing.T) {
	type row struct {
		Name    string `jrxml:"name"`
		Count   int    `json:"count,omitempty"`
		Hidden  string `jrxml:"-"`
		Plain   bool
		private int
	}
	for _, tc := range []struct {
		name   string
		source DataSource
		field  string
		want   any
	}{
		{"struct", FromSlice([]row{{Name: "a", Count: 1}, {Name: "b", Count: 2}, {Name: "c", Count: 3}}), "name", "a"},
		{"map", FromSlice([]map[string]any{{"name": "a"}, {"name": "b"}, {"name": "c"}}), "name", "a"},
		{"any", FromSlice([]any{[]any{"a"}, []any{"b"}, []any{"c"}}), "0", "a"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cur, err := tc.source.Open(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			defer cur.Close()
			for i := 0; i < 3; i++ {
				rec, ok, err := cur.Next(context.Background())
				if err != nil || !ok {
					t.Fatalf("Next %d: %v %v", i, ok, err)
				}
				if i == 0 {
					v, found := rec.Value(tc.field)
					if !found || v != tc.want {
						t.Fatalf("field %q = %v,%v", tc.field, v, found)
					}
					if v, found := rec.Value("missing"); found || v != nil {
						t.Fatalf("missing = %v,%v", v, found)
					}
				}
			}
			if rec, ok, err := cur.Next(context.Background()); rec != nil || ok || err != nil {
				t.Fatalf("exhausted: %v,%v,%v", rec, ok, err)
			}
		})
	}
}
func TestSliceStructTags(t *testing.T) {
	type row struct {
		Name    string `jrxml:"alias"`
		Hidden  int    `json:"-"`
		Plain   bool
		private int
	}
	src := FromSlice([]row{{Name: "x", Plain: true}})
	cur, _ := src.Open(context.Background())
	rec, _, _ := cur.Next(context.Background())
	for _, name := range []string{"Hidden", "private", "Name"} {
		if _, ok := rec.Value(name); ok {
			t.Errorf("unexpected %s", name)
		}
	}
	if v, ok := rec.Value("alias"); !ok || v != "x" {
		t.Errorf("alias = %v,%v", v, ok)
	}
	if v, ok := rec.Value("Plain"); !ok || v != true {
		t.Errorf("Plain = %v,%v", v, ok)
	}
}
func TestSliceCancellationAndClose(t *testing.T) {
	src := FromSlice([]map[string]int{{"n": 1}, {"n": 2}})
	ctx, cancel := context.WithCancel(context.Background())
	cur, _ := src.Open(ctx)
	_, _, _ = cur.Next(ctx)
	cancel()
	if _, _, err := cur.Next(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled Next: %v", err)
	}
	if err := cur.Close(); err != nil {
		t.Fatal(err)
	}
	if err := cur.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := src.Open(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled Open: %v", err)
	}
}
func TestSliceEmptyAndUnsupported(t *testing.T) {
	src := FromSlice([]any(nil))
	cur, _ := src.Open(context.Background())
	if _, ok, _ := cur.Next(context.Background()); ok {
		t.Fatal("nil slice yielded row")
	}
	bad := FromSlice([]any{42})
	cur, _ = bad.Open(context.Background())
	if _, _, err := cur.Next(context.Background()); err == nil {
		t.Fatal("scalar row accepted")
	}
}
