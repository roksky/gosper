package expression

import "time"

// Value is a value produced or consumed by an expression. The set of concrete
// value types is closed: implementations outside this package cannot create one.
//
// The concrete types are Null, Bool, Int, Decimal, String and Time. Decimal is
// currently a placeholder: its concrete representation is an open decision
// (OQ-1, feature EXP-008) and arithmetic on it is undefined until that lands.
type Value interface{ isValue() }

// Null is the absence of a value.
type Null struct{}

// Bool is a boolean value.
type Bool bool

// Int is a signed 64-bit integer value.
type Int int64

// Decimal is a placeholder for an arbitrary-precision decimal value. It
// currently carries the textual form only; do not rely on it for arithmetic
// until EXP-008 defines the representation.
type Decimal string

// String is a string value.
type String string

// Time is a date/time value. Locale and timezone semantics are deferred to
// EXP-011.
type Time time.Time

func (Null) isValue()    {}
func (Bool) isValue()    {}
func (Int) isValue()     {}
func (Decimal) isValue() {}
func (String) isValue()  {}
func (Time) isValue()    {}

// NullValue is the canonical null value.
var NullValue Value = Null{}

// Symbols resolves the three JRXML symbol namespaces for one evaluation scope.
//
// Implementations must be safe for concurrent use if shared. A missing name
// returns (nil, false) and must not panic. Names are case-sensitive.
type Symbols interface {
	Field(name string) (Value, bool)
	Parameter(name string) (Value, bool)
	Variable(name string) (Value, bool)
}

// Kind identifies a symbol namespace.
type Kind int

const (
	// KindField is the report field namespace ($F{...}).
	KindField Kind = iota
	// KindParameter is the report parameter namespace ($P{...}).
	KindParameter
	// KindVariable is the report variable namespace ($V{...}).
	KindVariable
)

// String returns the lowercase namespace name.
func (k Kind) String() string {
	switch k {
	case KindField:
		return "field"
	case KindParameter:
		return "parameter"
	case KindVariable:
		return "variable"
	default:
		return "unknown"
	}
}

// Lookup resolves name in s using the namespace kind. A nil s or an unknown
// kind returns (nil, false).
func (k Kind) Lookup(s Symbols, name string) (Value, bool) {
	if s == nil {
		return nil, false
	}
	switch k {
	case KindField:
		return s.Field(name)
	case KindParameter:
		return s.Parameter(name)
	case KindVariable:
		return s.Variable(name)
	default:
		return nil, false
	}
}

// Scope is an ordered list of Symbols used for nested evaluation. Lookups
// consult scopes from first to last and the first match wins, so inner scopes
// shadow outer ones. A nil entry is skipped.
type Scope []Symbols

// Lookup resolves name of the given kind against the scope chain.
func (s Scope) Lookup(kind Kind, name string) (Value, bool) {
	for _, syms := range s {
		if v, ok := kind.Lookup(syms, name); ok {
			return v, true
		}
	}
	return nil, false
}
