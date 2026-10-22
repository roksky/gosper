package report_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"image"
	"os"
	"strings"
	"testing"

	"github.com/roksky/gosper/report"
)

func compileMedulla(t *testing.T, name string) *report.CompiledReport {
	t.Helper()
	f, err := os.Open("../testdata/medulla/" + name + ".jrxml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	c, err := report.Compile(context.Background(), f)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func TestClassListTotalsAndLogo(t *testing.T) {
	c := compileMedulla(t, "students_class_list")
	rows := []map[string]any{{"name": "Alice", "gender": "female", "student_class": "One"}, {"name": "Bob", "gender": "male", "student_class": "One"}, {"name": "Cara", "gender": "female", "student_class": "Two"}}
	filled, err := report.Fill(context.Background(), c, report.Params{"logo": image.NewRGBA(image.Rect(0, 0, 2, 2))}, report.FromSlice(rows))
	if err != nil {
		t.Fatal(err)
	}
	var html bytes.Buffer
	if err := report.Export(context.Background(), filled, &html, report.FormatHTML); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Alice", "Bob", "Cara", "data:image/png;base64,", "#006699", "font-weight:bold"} {
		if !strings.Contains(html.String(), want) {
			t.Errorf("HTML missing %q", want)
		}
	}
	var out bytes.Buffer
	if err := report.Export(context.Background(), filled, &out, report.FormatCSV); err != nil {
		t.Fatal(err)
	}
	data, err := csv.NewReader(&out).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	// Final summary uses the final variable snapshot (3 students, 1 male, 2 female).
	text := ""
	for _, row := range data {
		text += strings.Join(row, "|") + "\n"
	}
	for _, want := range []string{"3", "1", "2"} {
		if !strings.Contains(text, want) {
			t.Errorf("summary missing %s: %s", want, text)
		}
	}
	var pdf bytes.Buffer
	if err := report.Export(context.Background(), filled, &pdf, report.FormatPDF); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(pdf.String(), "BI /W 2 /H 2") {
		t.Fatal("school logo missing from PDF")
	}
}
func TestGroupResetsAndExactDecimalTotals(t *testing.T) {
	c := compileMedulla(t, "finance_students_fee_balance")
	rows := []map[string]any{{"current_class": "One", "current_stream": "A", "invoiced_total": "0.10", "paid_total": "0.00", "balance": "0.10"}, {"current_class": "One", "current_stream": "A", "invoiced_total": "0.20", "paid_total": "0.10", "balance": "0.10"}, {"current_class": "Two", "current_stream": "B", "invoiced_total": "1.00", "paid_total": "0.50", "balance": "0.50"}}
	filled, err := report.Fill(context.Background(), c, nil, report.FromSlice(rows))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := report.Export(context.Background(), filled, &out, report.FormatCSV); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"0.3", "1.3", "0.6", "0.7"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing group/report total %s in %s", want, out.String())
		}
	}
}
func TestReceiptFormattingAndDefaults(t *testing.T) {
	c := compileMedulla(t, "system_receipt")
	filled, err := report.Fill(context.Background(), c, nil, report.FromSlice([]map[string]any{{"names": "Alice", "total": "1250.50", "fee_balance": "99.00"}}))
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range []report.Format{report.FormatHTML, report.FormatPDF, report.FormatXLSX, report.FormatCSV, report.FormatXML} {
		var out bytes.Buffer
		if err := report.Export(context.Background(), filled, &out, format); err != nil {
			t.Fatal(format, err)
		}
		if out.Len() == 0 {
			t.Fatal("empty receipt")
		}
	}
}
func TestPaginationUsesFinalReportPageNumber(t *testing.T) {
	c, err := report.Compile(context.Background(), strings.NewReader(`<jasperReport pageWidth="200" pageHeight="100" topMargin="5" bottomMargin="5"><field name="name"/><columnHeader><band height="10"/></columnHeader><detail><band height="40"><textField><reportElement width="100" height="20"/><textFieldExpression>$F{name}</textFieldExpression></textField></band></detail><pageFooter><band height="10"><textField><reportElement width="100" height="10"/><textFieldExpression>"Page "+$V{PAGE_NUMBER}</textFieldExpression></textField><textField evaluationTime="Report"><reportElement x="100" width="50" height="10"/><textFieldExpression>"of "+$V{PAGE_NUMBER}</textFieldExpression></textField></band></pageFooter></jasperReport>`))
	if err != nil {
		t.Fatal(err)
	}
	filled, err := report.Fill(context.Background(), c, nil, report.FromSlice([]map[string]string{{"name": "Alice"}, {"name": "Bob"}, {"name": "Cara"}}))
	if err != nil {
		t.Fatal(err)
	}
	if filled.PageCount() != 3 {
		t.Fatalf("pages=%d", filled.PageCount())
	}
	var out bytes.Buffer
	if err := report.Export(context.Background(), filled, &out, report.FormatCSV); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Page 1", "Page 2", "Page 3", "of 3"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing %s in %s", want, out.String())
		}
	}
}
