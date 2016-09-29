package util_test

import (
	"testing"

	. "github.com/doozr/qbot/util"
)

var stringPopTests = []struct {
	desc  string
	in    string
	first string
	rest  string
}{
	{"splits first token", "one two three", "one", "two three"},
	{"returns whole string", "onetwothree", "onetwothree", ""},
	{"trims whitespace", "one  two three  ", "one", "two three"},
	{"returns empty", "", "", ""},
}

func TestStringPop(t *testing.T) {
	for _, tt := range stringPopTests {
		first, rest := StringPop(tt.in)
		if first != tt.first {
			t.Errorf("It %s; expected '%s', received '%s'", tt.desc, tt.first, first)
		}
		if rest != tt.rest {
			t.Errorf("It %s; expected '%s', received '%s'", tt.desc, tt.rest, rest)
		}
	}
}
