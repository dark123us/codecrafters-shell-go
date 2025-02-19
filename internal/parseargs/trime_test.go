package parseargs

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
		{"'hello     script' 'shell''world'", []string{"hello     script", "shellworld"}},
		{"\"quz  hello\"  \"bar\"", []string{"quz  hello", "bar"}},
		{"\"bar\"  \"shell's\"  \"foo\"", []string{"bar", "shell's", "foo"}},
		{"\"before\\  after\"", []string{"before\\  after"}},
		{"world\\ \\ \\ \\ \\ \\ script", []string{"world      script"}},
		{"'shell\\\\\\nscript'", []string{"shell\\\\\\nscript"}},
		{"'example\\\"testhello\\\"shell'", []string{"example\\\"testhello\\\"shell"}},
	}

	for _, c := range cases {
		result := TrimString(c.arg)
		if !compareStringSlice(result, c.expected) {
			t.Errorf("TrimString(%q) = %v; expected %v", c.arg, result, c.expected)
		}
	}
}
