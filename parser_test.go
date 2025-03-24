package laks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	var tests = []struct{
		name string
		in []Token
		want []Expression
	}{
		{
			name: "empty",
			in: []Token{},
			want: nil,
		},
		{
			name: "literal",
			in: []Token{
				{T_INT, "44"},
				{T_SEMI, ";"},
			},
			want: []Expression{
				{T: E_LIT, Value: 44},
			},
		},	
		{
			name: "simple_add",
			in: []Token{
				{T_INT, "6"},
				{T_ADD, "+"},
				{T_INT, "7"},
				{T_SEMI, ";"},
			},
			want: []Expression{
				{
					T: E_OP, 
					Value: BO_ADD, 
					Left: &Expression{T: E_LIT, Value: 6},
					Right: &Expression{T: E_LIT, Value: 7},
				},
			},
		},
		{
			name: "prec1",
			in: []Token{
				{T_INT, "6"},
				{T_ADD, "+"},
				{T_INT, "7"},
				{T_MULT, "*"},
				{T_INT, "9"},
				{T_SEMI, ";"},
			},
			want: []Expression{
				{
					T: E_OP, 
					Value: BO_ADD, 
					Left: &Expression{T: E_LIT, Value: 6},
					Right: &Expression{
						T: E_OP,
						Value: BO_MULT,
						Left: &Expression{T: E_LIT, Value: 7},
						Right: &Expression{T: E_LIT, Value: 9},
					},
				},
			},
		},
		{
			name: "prec2",
			in: []Token{
				{T_INT, "6"},
				{T_MULT, "*"},
				{T_INT, "7"},
				{T_ADD, "+"},
				{T_INT, "9"},
				{T_SEMI, ";"},
			},
			want: []Expression{
				{
					T: E_OP, 
					Value: BO_ADD, 
					Left: &Expression{
						T: E_OP,
						Value: BO_MULT,
						Left: &Expression{T: E_LIT, Value: 6},
						Right: &Expression{T: E_LIT, Value: 7},
					},
					Right: &Expression{T: E_LIT, Value: 9},
				},
			},
		},	
	}

	for _, tst := range tests {
		t.Run(tst.name, func(tt *testing.T) {
			got, err := Parse(tst.in)
			if err != nil {
				tt.Fatalf("%s", err.Error())
			}
			if diff := cmp.Diff(tst.want, got); diff != "" {
				tt.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
