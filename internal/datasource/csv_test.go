package datasource

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
)

func TestCSVSource(t *testing.T) {
	f, err := os.Open("../../testdata/csv/basic.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	src := FromCSV(f)
	for run := 0; run < 2; run++ {
		cur, err := src.Open(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []struct{ name, city, note string }{{"Ada", "Nairobi", "hello, world"}, {"Biko", "Kampala", "line one\nline two"}} {
			rec, ok, err := cur.Next(context.Background())
			if err != nil || !ok {
				t.Fatalf("Next = %v,%v", ok, err)
			}
			for field, value := range map[string]string{"name": want.name, "city": want.city, "note": want.note} {
				if got, found := rec.Value(field); !found || got != value {
					t.Errorf("%s = %v,%v", field, got, found)
				}
			}
			if got, found := rec.Value("other"); found || got != nil {
				t.Errorf("missing = %v,%v", got, found)
			}
		}
		if rec, ok, err := cur.Next(context.Background()); rec != nil || ok || err != nil {
			t.Fatalf("end = %v,%v,%v", rec, ok, err)
		}
		cur.Close()
		if rec, ok, err := cur.Next(context.Background()); rec != nil || ok || err != nil {
			t.Fatalf("closed = %v,%v,%v", rec, ok, err)
		}
	}
}

func TestCSVDelimiterAndBOM(t *testing.T) {
	src := FromCSV(strings.NewReader("\ufeffname;amount\nAda;12.5\n"), WithCSVDelimiter(';'))
	cur, err := src.Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	rec, ok, err := cur.Next(context.Background())
	if err != nil || !ok {
		t.Fatalf("Next = %v,%v", ok, err)
	}
	if got, _ := rec.Value("name"); got != "Ada" {
		t.Fatal(got)
	}
	if got, _ := rec.Value("amount"); got != "12.5" {
		t.Fatal(got)
	}
}

func TestCSVInvalidInput(t *testing.T) {
	for _, input := range []string{"", "a,a\n1,2\n", "a,\n1,2\n", "a,b\n1\n", "a,b\n1,2,3\n", "a,b\n\"unfinished\n"} {
		_, err := FromCSV(strings.NewReader(input)).Open(context.Background())
		if d, ok := diag.As(err); !ok || d.Code != diag.CodeParse {
			t.Errorf("%q: %v", input, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := FromCSV(strings.NewReader("a\n1\n")).Open(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
