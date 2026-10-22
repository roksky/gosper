package jrxml

import (
	"fmt"
	"strings"

	"github.com/roksky/gosper/internal/diag"
)

// Normalize converts supported 6.x variants to one model shape. It mutates t
// and is idempotent. The source version must be inspected before normalization
// if a caller needs to retain its exact spelling.
func Normalize(t *Template) diag.List {
	var warnings diag.List
	if t == nil {
		warnings.Add(diag.New(diag.CodeValidation, diag.SeverityWarning, "cannot normalize nil template", nil))
		return warnings
	}
	if !t.Dialect.Known() {
		warnings.Add(diag.New(diag.CodeUnsupported, diag.SeverityWarning, "JRXML version is missing or unknown; 6.x defaults assumed", nil))
	} else if t.Dialect.Major != 6 {
		warnings.Add(diag.New(diag.CodeUnsupported, diag.SeverityWarning, fmt.Sprintf("JRXML version %s is outside the supported 6.x dialect", t.Dialect), nil))
	} else {
		t.Dialect = Dialect{Major: 6, Raw: "6"}
	}
	if strings.TrimSpace(t.Language) == "" {
		t.Language = "java"
	}
	for bi := range t.Bands {
		for ei := range t.Bands[bi].Elements {
			e := &t.Bands[bi].Elements[ei]
			if e.EvalTime == "" {
				e.EvalTime = "Now"
			}
			switch e.TextAdjust {
			case "", "StretchHeight", "CutText":
				if e.TextAdjust == "StretchHeight" {
					e.Stretch = true
				} else if e.TextAdjust == "CutText" {
					if e.Stretch {
						warnings.Add(diag.New(diag.CodeValidation, diag.SeverityWarning, "conflicting textAdjust and isStretchWithOverflow; textAdjust takes precedence", nil))
					}
					e.Stretch = false
				}
				e.TextAdjust = ""
			default:
				warnings.Add(diag.New(diag.CodeUnsupported, diag.SeverityWarning, fmt.Sprintf("uncertain textAdjust %q on %s", e.TextAdjust, e.Kind), nil))
			}
		}
	}
	return warnings
}
