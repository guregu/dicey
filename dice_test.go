package dicey

import (
	"testing"
)

type testCase struct {
	input string
	max   int
	min   int
	err   bool
}

var cases = []testCase{
	{
		input: "2d6",
		max:   12,
		min:   2,
	},
	{
		input: "10",
		max:   10,
		min:   10,
	},
	{
		input: "-10",
		max:   -10,
		min:   -10,
	},
	{
		input: "5d10+6d2",
		max:   62,
		min:   11,
	},
	{
		input: "5d10-6d2",
		max:   38,
		min:   -1,
	},
	{
		input: "1+2",
		max:   3,
		min:   3,
	},
	{
		input: "2+2d10",
		max:   22,
		min:   4,
	},
	{
		input: "2d10+2",
		max:   22,
		min:   4,
	},
	{
		input: "2-2d10",
		max:   -18,
		min:   0,
	},
	{
		input: "2d10-2",
		max:   18,
		min:   0,
	},
	{
		input: "6d9+4d20-2d6",
		max:   122,
		min:   8,
	},
}

func TestDice(t *testing.T) {
	for _, c := range cases {
		d, err := Parse(c.input)
		if err != nil && !c.err {
			t.Error("unexpected error for", c.input, c.err)
			continue
		}
		if c.err {
			if err == nil {
				t.Error("expected error for", c.input)
			}
			continue
		}
		if max := d.Max(); max != c.max {
			t.Error("bad max for", c.input, "expected:", c.max, "got:", max)
		}
		if min := d.Min(); min != c.min {
			t.Error("bad min for", c.input, "expected:", c.min, "got:", min)
		}
	}
}
