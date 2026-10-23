package expression

import (
	"context"
	"testing"
	"time"
)

func TestJasperExpressions(t *testing.T) {
	scope := Scope{fakeSymbols{fields: map[string]Value{"gender": String(" FEMALE "), "count": Int(3), "amount": Decimal("0.2")}, parameters: map[string]Value{"REPORT_TIME": Time(time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC))}}}
	for _, tc := range []struct {
		source string
		want   Value
	}{
		{`$F{gender}.trim().toLowerCase().equals("female") ? 1 : 0`, Int(1)},
		{`true ? "yes" : 1/0`, String("yes")},
		{`"Page "+$F{count}`, String("Page 3")},
		{`"new Date()"`, String("new Date()")},
		{`0.1 + $F{amount}`, Decimal("0.3")},
		{`0.000000000000000001 + 0.000000000000000002`, Decimal("0.000000000000000003")},
		{`1000000000000000000.0 * 1000`, Decimal("1000000000000000000000")},
		{`new DecimalFormat("#,##0.00").format(1250.5)`, String("1,250.50")},
		{`new SimpleDateFormat("MMM d, yyyy").format(new Date())`, String("Oct 22, 2026")},
	} {
		p, err := Compile(tc.source)
		if err != nil {
			t.Fatal(tc.source, err)
		}
		v, err := p.Eval(context.Background(), scope)
		if err != nil || v != tc.want {
			t.Errorf("%s: %v %v want %v", tc.source, v, err, tc.want)
		}
	}
}

func TestMissingGenderMethodChainContributesZero(t *testing.T) {
	p, err := Compile(`(($F{gender}.trim().toLowerCase().equals("male")) ? 1 : 0)`)
	if err != nil {
		t.Fatal(err)
	}
	for _, fields := range []map[string]Value{nil, {"gender": NullValue}, {"gender": String("")}, {"gender": String("OTHER")}} {
		value, err := p.Eval(context.Background(), Scope{fakeSymbols{fields: fields}})
		if err != nil || value != Int(0) {
			t.Fatalf("fields=%v result=%v error=%v", fields, value, err)
		}
	}
}
