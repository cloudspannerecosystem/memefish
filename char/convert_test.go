package char

import (
	"testing"
)

func TestToUpper(t *testing.T) {
	for _, tc := range []struct {
		input, want string
	}{
		{"ABC", "ABC"},
		{"abc", "ABC"},
		{"-=-", "-=-"},
		{"αβγ", "αβγ"},
	} {
		got := ToUpper(tc.input)
		if tc.want != got {
			t.Errorf("ToUpper(%q): %q (want) != %q (got)", tc.input, tc.want, got)
		}
	}
}

func TestEqualFold(t *testing.T) {
	for _, tc := range []struct {
		s, t string
		want bool
	}{
		{"ABC", "ABC", true},
		{"abc", "ABC", true},
		{"ABC", "abc", true},
		{"-=-", "-=-", true},
		{"αβγ", "ΑΒΓ", false},
	} {
		got := EqualFold(tc.s, tc.t)
		if tc.want != got {
			t.Errorf("EqualFold(%q, %q) != %t", tc.s, tc.t, tc.want)
		}
	}
}
