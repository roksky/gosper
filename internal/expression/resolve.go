package expression

import (
	"fmt"
	"strings"
)

// Reference identifies a symbol. Explicit $F/$P/$V references always use
// their named namespace; bare names use variable, parameter, field order.
type Reference struct {
	Kind     Kind
	Name     string
	Explicit bool
}

// ParseReference parses a single reference. Whitespace around a braced name
// is ignored; empty names, nested braces and malformed references are errors.
func ParseReference(src string) (Reference, error) {
	if len(src) >= 3 && src[0] == '$' {
		var kind Kind
		switch src[1] {
		case 'F':
			kind = KindField
		case 'P':
			kind = KindParameter
		case 'V':
			kind = KindVariable
		default:
			return Reference{}, fmt.Errorf("unknown reference kind in %q", src)
		}
		if src[2] != '{' || src[len(src)-1] != '}' {
			return Reference{}, fmt.Errorf("malformed reference %q", src)
		}
		name := strings.TrimSpace(src[3 : len(src)-1])
		if !validName(name) {
			return Reference{}, fmt.Errorf("invalid reference name in %q", src)
		}
		return Reference{Kind: kind, Name: name, Explicit: true}, nil
	}
	if !validName(src) {
		return Reference{}, fmt.Errorf("invalid reference %q", src)
	}
	return Reference{Name: src}, nil
}

func validName(name string) bool {
	if name == "" || strings.ContainsAny(name, "{}$ \t\r\n") {
		return false
	}
	return true
}

// Resolve looks up a reference in the supplied scope. Missing names return
// (nil, false), including when the scope is empty. An explicit namespace never
// falls back to another namespace. For bare names, variable wins over
// parameter, which wins over field; within each namespace inner scope wins.
func (r Reference) Resolve(scope Scope) (Value, bool) {
	if r.Explicit {
		return scope.Lookup(r.Kind, r.Name)
	}
	for _, kind := range [...]Kind{KindVariable, KindParameter, KindField} {
		if v, ok := scope.Lookup(kind, r.Name); ok {
			return v, true
		}
	}
	return nil, false
}
