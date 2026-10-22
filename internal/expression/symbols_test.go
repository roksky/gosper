package expression

import (
	"testing"
	"time"
)

// fakeSymbols is a map-backed Symbols implementation used to exercise the
// contract.
type fakeSymbols struct {
	fields     map[string]Value
	parameters map[string]Value
	variables  map[string]Value
}

func (f fakeSymbols) Field(name string) (Value, bool)     { return lookup(f.fields, name) }
func (f fakeSymbols) Parameter(name string) (Value, bool) { return lookup(f.parameters, name) }
func (f fakeSymbols) Variable(name string) (Value, bool)  { return lookup(f.variables, name) }

func lookup(m map[string]Value, name string) (Value, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m[name]
	return v, ok
}

var _ Symbols = fakeSymbols{}

func TestValueTypes(t *testing.T) {
	values := []Value{
		Null{},
		Bool(true),
		Int(42),
		Decimal("1.50"),
		String("x"),
		Time(time.Unix(0, 0).UTC()),
	}
	for i, v := range values {
		if v == nil {
			t.Fatalf("value %d is nil", i)
		}
	}
	if NullValue == nil {
		t.Fatal("NullValue must not be nil")
	}
}

func TestKindLookup(t *testing.T) {
	s := fakeSymbols{
		fields:     map[string]Value{"f": Int(1)},
		parameters: map[string]Value{"p": String("hello")},
		variables:  map[string]Value{"v": Bool(true)},
	}

	cases := []struct {
		kind Kind
		name string
		want Value
		ok   bool
	}{
		{KindField, "f", Int(1), true},
		{KindParameter, "p", String("hello"), true},
		{KindVariable, "v", Bool(true), true},
		{KindField, "p", nil, false}, // namespaces are separate
		{KindParameter, "missing", nil, false},
	}
	for _, tc := range cases {
		got, ok := tc.kind.Lookup(s, tc.name)
		if ok != tc.ok || got != tc.want {
			t.Fatalf("%s.Lookup(%q) = (%v, %v), want (%v, %v)", tc.kind, tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestLookupNilSymbolsAndUnknownKind(t *testing.T) {
	if _, ok := KindField.Lookup(nil, "x"); ok {
		t.Fatal("nil Symbols must miss")
	}
	if _, ok := Kind(99).Lookup(fakeSymbols{}, "x"); ok {
		t.Fatal("unknown kind must miss")
	}
}

func TestScopePrecedence(t *testing.T) {
	outer := fakeSymbols{variables: map[string]Value{"x": Int(1), "n": Int(10)}}
	inner := fakeSymbols{variables: map[string]Value{"x": Int(2)}}

	// Inner (first) shadows outer (second).
	scope := Scope{inner, outer}
	if v, ok := scope.Lookup(KindVariable, "x"); !ok || v != Int(2) {
		t.Fatalf("scope x = (%v, %v), want inner value 2", v, ok)
	}
	if v, ok := scope.Lookup(KindVariable, "n"); !ok || v != Int(10) {
		t.Fatalf("scope n = (%v, %v), want outer value 10", v, ok)
	}
	if _, ok := scope.Lookup(KindVariable, "missing"); ok {
		t.Fatal("missing name must miss")
	}
}

func TestScopeSkipsNil(t *testing.T) {
	outer := fakeSymbols{fields: map[string]Value{"a": String("a")}}
	scope := Scope{nil, outer}
	if v, ok := scope.Lookup(KindField, "a"); !ok || v != String("a") {
		t.Fatalf("scope a = (%v, %v), want a", v, ok)
	}
	if _, ok := (Scope{}).Lookup(KindField, "a"); ok {
		t.Fatal("empty scope must miss")
	}
}

func TestKindString(t *testing.T) {
	for _, tc := range []struct {
		kind Kind
		want string
	}{{KindField, "field"}, {KindParameter, "parameter"}, {KindVariable, "variable"}, {Kind(9), "unknown"}} {
		if got := tc.kind.String(); got != tc.want {
			t.Fatalf("Kind(%d).String() = %q, want %q", tc.kind, got, tc.want)
		}
	}
}
