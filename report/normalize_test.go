package report

import (
	"context"
	"os"
	"testing"
)

func TestCompileSixXVariants(t *testing.T) {
	for _, name := range []string{"JRX-020-001-legacy.jrxml", "JRX-020-001-modern.jrxml"} {
		f, err := os.Open("../testdata/jrxml/versions/" + name)
		if err != nil {
			t.Fatal(err)
		}
		compiled, err := Compile(context.Background(), f)
		f.Close()
		if err != nil || compiled == nil {
			t.Errorf("Compile(%s) = %v, %v", name, compiled, err)
		}
	}
}
