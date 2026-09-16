package money

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
		{"-$4.50", -450},
		{"$-4.50", -450},
		{"-$1,234.56", -123456},
		{"+$3", 300},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		require.NoError(t, err, "Parse(%q)", c.in)
		require.Equal(t, c.want, got, "Parse(%q)", c.in)
	}
}

func TestParseInvalid(t *testing.T) {
	for _, in := range []string{"", "abc", "1.234", "1.2.3", "$", "-", "-$", "--1", " "} {
		_, err := Parse(in)
		require.Error(t, err, "Parse(%q)", in)
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
		require.Equal(t, c.want, Format(c.in), "Format(%d)", c.in)
	}
}

func TestParseIntoFormat(t *testing.T) {
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
		require.NoError(t, err, "Parse(%q)", c.in)

		formatted := Format(cents)
		require.Equal(t, c.want, formatted, "round trip %q", c.in)

		// Format's output must parse back to the same cents.
		back, err := Parse(formatted)
		require.NoError(t, err, "Parse(Format(%d))", cents)
		require.Equal(t, cents, back, "Parse(Format(%d))", cents)
	}
}
