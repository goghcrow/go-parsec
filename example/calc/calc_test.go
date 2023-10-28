package calc

import (
	"testing"
)

func TestRec(t *testing.T) {
	for _, tt := range []struct {
		input   string
		expect  float64
		expectS string
	}{
		{"1", 1, "1"},
		{"+1.5", 1.5, "(+ 1.5)"},
		{"-0.5", -0.5, "(- 0.5)"},
		{"1 + 2", 3, "(+ 1 2)"},
		{"1 - 2", -1, "(- 1 2)"},
		{"1 * 2", 2, "(* 1 2)"},
		{"1 / 2", 0.5, "(/ 1 2)"},
		{"1 + 2 * 3 + 4", 11, "(+ (+ 1 (* 2 3)) 4)"},
		{"1 + 2 + 3", 6, "(+ (+ 1 2) 3)"},
		{"(1 + 2) * (3 + 4)", 21, "(* (+ 1 2) (+ 3 4))"},
		{"1.2--3.4", 4.6, "(- 1.2 (- 3.4))"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			v := Calc(tt.input)
			if tt.expect != v {
				t.Errorf("expect %f actual %f", tt.expect, v)
			}
			s := Show(tt.input)
			if tt.expectS != s {
				t.Errorf("expect %s actual %s", tt.expectS, s)
			}
		})
	}
}
