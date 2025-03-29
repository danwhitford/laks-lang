package laks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		name string
		in   []Token
		want []Expression
	}{
		{
			name: "empty",
			in:   []Token{},
			want: nil,
		},
		{
			name: "literal",
			in: []Token{
				{T_INT, "44"},
				{T_SEMI, ";"},
			},
			want: []Expression{
				{T: E_LIT, Value: int64(44)},
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
					T:     E_BINOP,
					Value: BinaryExpression{BO_ADD, Expression{T: E_LIT, Value: int64(6)}, Expression{T: E_LIT, Value: int64(7)}},
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
					T: E_BINOP,
					Value: BinaryExpression{
						BO_ADD,
						Expression{T: E_LIT, Value: int64(6)},
						Expression{
							T: E_BINOP,
							Value: BinaryExpression{
								BO_MULT,
								Expression{T: E_LIT, Value: int64(7)},
								Expression{T: E_LIT, Value: int64(9)},
							},
						},
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
					T: E_BINOP,
					Value: BinaryExpression{
						BO_ADD,
						Expression{
							T: E_BINOP,
							Value: BinaryExpression{
								BO_MULT,
								Expression{T: E_LIT, Value: int64(6)},
								Expression{T: E_LIT, Value: int64(7)},
							},
						},
						Expression{T: E_LIT, Value: int64(9)},
					},
				},
			},
		},
		{
			name: "print something",
			in: []Token{
				{T_KEYWORD, "print"},
				{T_INT, "7"},
				{T_MULT, "*"},
				{T_INT, "8"},
				{T_SEMI, ";"},
			},
			want: []Expression{
				{
					T: E_PRINT,
					Value: Expression{
						T: E_BINOP,
						Value: BinaryExpression{
							BO_MULT,
							Expression{E_LIT, int64(7)},
							Expression{E_LIT, int64(8)},
						},
					},
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
