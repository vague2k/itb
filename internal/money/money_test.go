package money

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"12", 1200},
		{"12.34", 1234},
		{"$1,234.56", 123456},
		{"0.05", 5},
		{".5", 50},
		{"-4.5", -450},
		{"+3", 300},
		{"  7.00  ", 700},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if err != nil {
			t.Fatalf("Parse(%q) unexpected error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("Parse(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	for _, in := range []string{"", "abc", "1.234", "1.2.3", "$", "--1"} {
		if _, err := Parse(in); err == nil {
			t.Errorf("Parse(%q) expected error", in)
		}
	}
}

func TestFormat(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "$0.00"},
		{5, "$0.05"},
		{1200, "$12.00"},
		{123456, "$1,234.56"},
		{100000000, "$1,000,000.00"},
		{-450, "-$4.50"},
	}
	for _, c := range cases {
		if got := Format(c.in); got != c.want {
			t.Errorf("Format(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"12.34", "$12.34"},
		{"0.05", "$0.05"},
		{"-4.50", "-$4.50"},
		{"999999.99", "$999,999.99"},
	}
	for _, c := range cases {
		cents, err := Parse(c.in)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.in, err)
		}
		if got := Format(cents); got != c.want {
			t.Errorf("round trip %q = %q, want %q", c.in, got, c.want)
		}
	}
}
