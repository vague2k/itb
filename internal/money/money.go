package money

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalid is returned when a user-entered amount cannot be parsed.
var ErrInvalid = errors.New("invalid amount")

// Parse converts a user-entered amount such as "$1,234.56", "12", or "-4.5"
// into an integer number of cents.
func Parse(s string) (int64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return 0, ErrInvalid
	}

	negative := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	s = strings.TrimPrefix(s, "+")

	whole, frac, _ := strings.Cut(s, ".")
	if whole == "" {
		whole = "0"
	}

	dollars, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || dollars < 0 {
		return 0, ErrInvalid
	}
	if len(frac) > 2 {
		return 0, ErrInvalid
	}
	for len(frac) < 2 {
		frac += "0"
	}
	cents, err := strconv.ParseInt(frac, 10, 64)
	if err != nil || cents < 0 {
		return 0, ErrInvalid
	}

	total := dollars*100 + cents
	if negative {
		return -total, nil
	}
	return total, nil
}

// Format renders an integer number of cents as "$1,234.56", or "-$12.00" when
// negative.
func Format(cents int64) string {
	if cents < 0 {
		return fmt.Sprintf("-$%s.%02d", groupThousands(strconv.FormatInt(-cents/100, 10)), -cents%100)
	}
	return fmt.Sprintf("$%s.%02d", groupThousands(strconv.FormatInt(cents/100, 10)), cents%100)
}

func groupThousands(digits string) string {
	if len(digits) <= 3 {
		return digits
	}
	var b strings.Builder
	if pre := len(digits) % 3; pre > 0 {
		b.WriteString(digits[:pre])
	}
	for i := len(digits) % 3; i < len(digits); i += 3 {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		b.WriteString(digits[i : i+3])
	}
	return b.String()
}
