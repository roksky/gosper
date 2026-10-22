package expression

import "testing"

func TestReferenceResolution(t *testing.T) {
	scope := Scope{fakeSymbols{
		fields:     map[string]Value{"x": Int(1), "field": Int(4)},
		parameters: map[string]Value{"x": Int(2), "parameter": Int(5)},
		variables:  map[string]Value{"x": Int(3), "variable": Int(6)},
	}}
	for _, tc := range []struct {
		src  string
		want Value
		ok   bool
	}{
		{"$F{x}", Int(1), true}, {"$P{x}", Int(2), true},
		{"$V{x}", Int(3), true}, {"x", Int(3), true},
		{"field", Int(4), true}, {"parameter", Int(5), true},
		{"variable", Int(6), true}, {"$F{ field }", Int(4), true},
		{"$V{missing}", nil, false}, {"missing", nil, false},
	} {
		r, err := ParseReference(tc.src)
		if err != nil {
			t.Fatalf("ParseReference(%q): %v", tc.src, err)
		}
		got, ok := r.Resolve(scope)
		if got != tc.want || ok != tc.ok {
			t.Errorf("%q: got (%v, %v), want (%v, %v)", tc.src, got, ok, tc.want, tc.ok)
		}
	}
}

func TestReferenceScopeAndNull(t *testing.T) {
	scope := Scope{
		fakeSymbols{variables: map[string]Value{"x": NullValue}},
		fakeSymbols{variables: map[string]Value{"x": Int(1)}},
	}
	r, _ := ParseReference("$V{x}")
	if got, ok := r.Resolve(scope); !ok || got != NullValue {
		t.Fatalf("shadowed null = (%v, %v)", got, ok)
	}
}

func TestInvalidReferences(t *testing.T) {
	for _, src := range []string{"", "$F{}", "$F{ }", "$F{x{y}}", "$F{x}tail", "$X{x}", "$F{x", "x y", "$"} {
		if _, err := ParseReference(src); err == nil {
			t.Errorf("ParseReference(%q) unexpectedly succeeded", src)
		}
	}
}
