package report_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/roksky/gosper/report"
)

func TestJRXMLQueryDefaults(t *testing.T) {
	const source = `<jasperReport name="sql"><parameter name="school_id"/><parameter name="label"><defaultValueExpression>"default"</defaultValueExpression></parameter><queryString language="SQL"><![CDATA[select $P{school_id}, $P{label}, $P{school_id}]]></queryString><detail><band height="10"/></detail></jasperReport>`
	c, err := report.Compile(context.Background(), strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	q, args, err := c.Query(context.Background(), report.Params{"school_id": "tenant"}, func(i int) string { return fmt.Sprintf("$%d", i) })
	if err != nil || q != "select $1, $2, $3" || !reflect.DeepEqual(args, []any{"tenant", "default", "tenant"}) {
		t.Fatalf("query=%q args=%v err=%v", q, args, err)
	}
}

func TestBindQuery(t *testing.T) {
	q, args, err := report.BindQuery("select '$P{x}', \"$P{x}\", $P{x} -- $P{x}\n/* $P{x} */", report.Params{"x": "'; DROP TABLE report; --"}, nil)
	if err != nil || q != "select '$P{x}', \"$P{x}\", ? -- $P{x}\n/* $P{x} */" || len(args) != 1 || args[0] != "'; DROP TABLE report; --" {
		t.Fatalf("%q %v %v", q, args, err)
	}
	for _, q := range []string{"select $P{missing}", "select $P!{x}", "select $X{IN,id,x}", "select $P{x", "select 'oops", "/* open", "select $$open"} {
		if _, _, err := report.BindQuery(q, nil, nil); err == nil {
			t.Errorf("accepted %q", q)
		}
	}
}

func TestRealExports(t *testing.T) {
	filled, err := report.Fill(context.Background(), fixture(t), nil, report.FromSlice([]map[string]string{{"name": "<Alice>&"}, {"name": "Bob"}, {"name": "Cara"}}))
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range []report.Format{report.FormatPDF, report.FormatHTML, report.FormatCSV, report.FormatXML, report.FormatXLSX} {
		t.Run(format.String(), func(t *testing.T) {
			var out bytes.Buffer
			if err := report.Export(context.Background(), filled, &out, format); err != nil {
				t.Fatal(err)
			}
			switch format {
			case report.FormatPDF:
				if !strings.HasPrefix(out.String(), "%PDF-") {
					t.Fatal("invalid PDF")
				}
			case report.FormatHTML:
				if !strings.Contains(out.String(), "&lt;Alice&gt;&amp;!") || strings.Contains(out.String(), "<Alice>") {
					t.Fatal(out.String())
				}
			case report.FormatCSV:
				rows, err := csv.NewReader(&out).ReadAll()
				if err != nil || len(rows) != 3 || rows[0][0] != "<Alice>&!" {
					t.Fatalf("%v %v", rows, err)
				}
			case report.FormatXML:
				var doc any
				if err := xml.Unmarshal(out.Bytes(), &doc); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(out.String(), "<page>") {
					t.Fatal(out.String())
				}
			case report.FormatXLSX:
				z, err := zip.NewReader(bytes.NewReader(out.Bytes()), int64(out.Len()))
				if err != nil {
					t.Fatal(err)
				}
				var sheet string
				for _, f := range z.File {
					r, err := f.Open()
					if err != nil {
						t.Fatal(err)
					}
					b, err := io.ReadAll(r)
					r.Close()
					if err != nil {
						t.Fatal(err)
					}
					var doc any
					if err := xml.Unmarshal(b, &doc); err != nil {
						t.Fatalf("%s: %v", f.Name, err)
					}
					if f.Name == "xl/worksheets/sheet1.xml" {
						sheet = string(b)
					}
				}
				if !strings.Contains(sheet, "&lt;Alice&gt;&amp;!") {
					t.Fatal(sheet)
				}
			}
		})
	}
}

func TestCompileRetainsUnsupportedFeatureErrors(t *testing.T) {
	if _, err := report.Compile(context.Background(), strings.NewReader(`<jasperReport><chart/></jasperReport>`)); err == nil {
		t.Fatal("unsupported image silently ignored")
	}
}

func TestDollarQuotedSQLPreservesLiteralParameters(t *testing.T) {
	q, args, err := report.BindQuery(`select $$ $P{literal} $$, $fn$ $P{literal} $fn$, $P{value}`, report.Params{"value": 42}, nil)
	if err != nil || q != `select $$ $P{literal} $$, $fn$ $P{literal} $fn$, ?` || len(args) != 1 {
		t.Fatalf("%q %v %v", q, args, err)
	}
}

func TestParameterNamesAreIndependent(t *testing.T) {
	c := fixture(t)
	names := c.ParameterNames()
	if len(names) != 1 || names[0] != "suffix" {
		t.Fatalf("unexpected names: %v", names)
	}
	names[0] = "changed"
	if c.ParameterNames()[0] != "suffix" {
		t.Fatal("parameter names mutate compiled report")
	}
}
