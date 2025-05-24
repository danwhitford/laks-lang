package laks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		name string
		in   []Token
		want []Statement
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
			want: []Statement{
				LiteralExpression{Value{VAL_INT, int64(44)}},
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
			want: []Statement{
				BinaryExpression{
					BO_ADD,
					LiteralExpression{Value{VAL_INT, int64(6)}},
					LiteralExpression{Value{VAL_INT, int64(7)}},
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
			want: []Statement{
				BinaryExpression{
					BO_ADD,
					LiteralExpression{Value{VAL_INT, int64(6)}},
					BinaryExpression{
						BO_MULT,
						LiteralExpression{Value{VAL_INT, int64(7)}},
						LiteralExpression{Value{VAL_INT, int64(9)}},
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
			want: []Statement{
				BinaryExpression{
					BO_ADD,
					BinaryExpression{
						BO_MULT,
						LiteralExpression{Value{VAL_INT, int64(6)}},
						LiteralExpression{Value{VAL_INT, int64(7)}},
					},
					LiteralExpression{Value{VAL_INT, int64(9)}},
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
			want: []Statement{
				PrintStatment{
					BinaryExpression{
						BO_MULT,
						LiteralExpression{Value{VAL_INT, int64(7)}},
						LiteralExpression{Value{VAL_INT, int64(8)}},
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

func TestErorrs(t *testing.T) {
	in := []Token{
		{T_KEYWORD, "print"},
		{T_STRING, "foobar!"},
	}
	_, err := Parse(in)
	if err == nil {
		t.Fatalf("wanted error")
	}
}
