package datasource

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
)

func TestJSONSource(t *testing.T) {
	f, err := os.Open("../../testdata/json/basic.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	src := FromJSON(f)
	for run := 0; run < 2; run++ {
		cur, err := src.Open(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		for i, want := range []string{"Ada", "Biko"} {
			rec, ok, err := cur.Next(context.Background())
			if err != nil || !ok {
				t.Fatalf("row %d: %v, %v", i, ok, err)
			}
			if got, _ := rec.Value("name"); got != want {
				t.Fatalf("name = %v", got)
			}
			if i == 0 {
				for key, want := range map[string]any{"count": float64(2), "active": true, "missing": nil} {
					if got, found := rec.Value(key); !found || got != want {
						t.Errorf("%s = %v,%v", key, got, found)
					}
				}
				if got, _ := rec.Value("meta"); got.(map[string]any)["city"] != "Nairobi" {
					t.Errorf("meta = %v", got)
				}
				if got, _ := rec.Value("tags"); len(got.([]any)) != 2 {
					t.Errorf("tags = %v", got)
				}
			} else if got, _ := rec.Value("count"); got != 3.5 {
				t.Errorf("count = %v", got)
			}
			if got, found := rec.Value("unknown"); found || got != nil {
				t.Errorf("missing = %v,%v", got, found)
			}
		}
		if rec, ok, err := cur.Next(context.Background()); rec != nil || ok || err != nil {
			t.Fatalf("end = %v,%v,%v", rec, ok, err)
		}
		if err := cur.Close(); err != nil {
			t.Fatal(err)
		}
		if rec, ok, err := cur.Next(context.Background()); rec != nil || ok || err != nil {
			t.Fatalf("closed = %v,%v,%v", rec, ok, err)
		}
	}
}

func TestJSONInvalidInput(t *testing.T) {
	for _, input := range []string{`{}`, `null`, `[1]`, `[null]`, `[{}] true`, `[{`, ``} {
		_, err := FromJSON(strings.NewReader(input)).Open(context.Background())
		if d, ok := diag.As(err); !ok || d.Code != diag.CodeParse {
			t.Errorf("%q: %v", input, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := FromJSON(strings.NewReader(`[]`)).Open(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
