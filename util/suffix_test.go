package util_test

import (
	"testing"

	. "github.com/doozr/qbot/util"
)

var suffixTests = []struct {
	num    int
	suffix string
}{
	{1, "st"},
	{2, "nd"},
	{3, "rd"},
	{4, "th"},
	{5, "th"},
	{6, "th"},
	{7, "th"},
	{8, "th"},
	{9, "th"},
	{0, "th"},
	{6511, "th"},
	{6512, "th"},
	{6513, "th"},
	{6541, "st"},
	{6542, "nd"},
	{6543, "rd"},
	{6544, "th"},
	{6545, "th"},
	{6546, "th"},
	{6547, "th"},
	{6548, "th"},
	{6549, "th"},
	{6540, "th"},
}

func TestSuffix(t *testing.T) {
	for _, tt := range suffixTests {
		suffix := Suffix(tt.num)
		if suffix != tt.suffix {
			t.Errorf("Expected suffix '%s' for number '%d', received '%s'", tt.suffix, tt.num, suffix)
		}
	}
}
