package laks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenise(t *testing.T) {
	var tests = []struct {
		in   string
		want []Token
	}{
		{
			in:   "",
			want: nil,
		},
		{
			in: "4",
			want: []Token{
				{T_INT, "4"},
			},
		},
		{
			in: "4;",
			want: []Token{
				{T_INT, "4"},
				{T_SEMI, ";"},
			},
		},
		{
			in: "8 * 7;",
			want: []Token{
				{T_INT, "8"},
				{T_MULT, "*"},
				{T_INT, "7"},
				{T_SEMI, ";"},
			},
		},
		{
			in: "\n\n2+2\t\t;",
			want: []Token{
				{T_INT, "2"},
				{T_ADD, "+"},
				{T_INT, "2"},
				{T_SEMI, ";"},
			},
		},
		{
			in: "+/-*",
			want: []Token{
				{T_ADD, "+"},
				{T_DIV, "/"},
				{T_MINUS, "-"},
				{T_MULT, "*"},
			},
		},
		{
			in: "print 7*8;",
			want: []Token{
				{T_KEYWORD, "print"},
				{T_INT, "7"},
				{T_MULT, "*"},
				{T_INT, "8"},
				{T_SEMI, ";"},
			},
		},
		{
			in: "print 7*8; # this is a comment",
			want: []Token{
				{T_KEYWORD, "print"},
				{T_INT, "7"},
				{T_MULT, "*"},
				{T_INT, "8"},
				{T_SEMI, ";"},
			},
		},
		{
			in: "# this is a comment\nprint 7*8;",
			want: []Token{
				{T_KEYWORD, "print"},
				{T_INT, "7"},
				{T_MULT, "*"},
				{T_INT, "8"},
				{T_SEMI, ";"},
			},
		},
		{
			in: "==",
			want: []Token{
				{T_EQ_EQ, "=="},
			},
		},
		{
			in: "true == false",
			want: []Token{
				{T_KEYWORD, "true"},
				{T_EQ_EQ, "=="},
				{T_KEYWORD, "false"},
			},
		},
		{
			in: "print \"foobar!\"",
			want: []Token{
				{T_KEYWORD, "print"},
				{T_STRING, "foobar!"},
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.in, func(tt *testing.T) {
			got, err := Tokenise([]byte(tst.in))
			if err != nil {
				tt.Fatalf("%s", err.Error())
			}
			if diff := cmp.Diff(tst.want, got); diff != "" {
				tt.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
