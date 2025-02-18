package main

import (
	"testing"
)

func compareStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestTrimString(t *testing.T) {
	cases := []struct {
		arg      string
		expected []string
	}{
		{"echo hello  world", []string{"echo", "hello", "world"}},
		{"echo 'hello  world'", []string{"echo", "hello  world"}},
		{"'echo ' ' hello world '", []string{"echo ", " hello world "}},
	}

	for _, c := range cases {
		result := trimString(c.arg)
		if !compareStringSlice(result, c.expected) {
			t.Errorf("trimString(%q) = %v; expected %v", c.arg, result, c.expected)
		}
	}
}
