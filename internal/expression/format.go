package expression

import (
	"fmt"
	"strings"
	"time"
)

// FormatValue handles the date and decimal patterns used by Jasper templates.
func FormatValue(value Value, pattern string) (string, error) {
	if isNull(value) {
		return "", nil
	}
	if date, ok := value.(Time); ok {
		replacer := strings.NewReplacer("EEEEE", "Monday", "EEEE", "Monday", "EEE", "Mon", "EE", "Mon", "MMMMM", "January", "MMMM", "January", "MMM", "Jan", "MM", "01", "yyyy", "2006", "yy", "06", "dd", "02", "HH", "15", "mm", "04", "ss", "05", "d", "2", "M", "1")
		return time.Time(date).Format(replacer.Replace(pattern)), nil
	}
	number, ok := numberRat(value)
	if !ok {
		return scalarString(value), nil
	}
	decimals := 0
	optional := false
	if dot := strings.IndexByte(pattern, '.'); dot >= 0 {
		decimals = len(pattern) - dot - 1
		optional = strings.Contains(pattern[dot+1:], "#")
	}
	if decimals > 16 {
		return "", fmt.Errorf("numeric pattern exceeds 16 decimal places")
	}
	formatted := number.FloatString(decimals)
	if optional {
		formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	}
	if strings.Contains(pattern, ",") {
		pieces := strings.SplitN(formatted, ".", 2)
		digits := pieces[0]
		negative := strings.HasPrefix(digits, "-")
		digits = strings.TrimPrefix(digits, "-")
		for i := len(digits) - 3; i > 0; i -= 3 {
			digits = digits[:i] + "," + digits[i:]
		}
		if negative {
			digits = "-" + digits
		}
		formatted = digits
		if len(pieces) > 1 {
			formatted += "." + pieces[1]
		}
	}
	return formatted, nil
}
