package laks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompile(t *testing.T) {
	var tests = []struct {
		name string
		in   []Statement
		want []byte
	}{
		{
			name: "literal",
			in: []Statement{
				LiteralExpression{IntValue(int64(14))},
			},
			want: []byte{
				byte(OP_PUSH),
				byte(VAL_INT),
				14, 0, 0, 0, 0, 0, 0, 0, // 14
			},
		},
		{
			name: "expradd",
			in: []Statement{
				BinaryExpression{
					Op:    BO_ADD,
					Left:  LiteralExpression{IntValue(int64(7))},
					Right: LiteralExpression{IntValue(int64(9))},
				},
			},
			want: []byte{
				byte(OP_PUSH),
				byte(VAL_INT),
				7, 0, 0, 0, 0, 0, 0, 0, // 7
				byte(OP_PUSH),
				byte(VAL_INT),
				9, 0, 0, 0, 0, 0, 0, 0, // 9
				byte(OP_ADD),
			},
		},
		{
			name: "print expression",
			in: []Statement{
				PrintStatment{
					Expr: BinaryExpression{
						Op:    BO_MULT,
						Left:  LiteralExpression{IntValue(int64(7))},
						Right: LiteralExpression{IntValue(int64(9))},
					},
				},
			},
			want: []byte{
				byte(OP_PUSH),
				byte(VAL_INT),
				7, 0, 0, 0, 0, 0, 0, 0, // 7
				byte(VAL_INT),
				byte(OP_PUSH),
				9, 0, 0, 0, 0, 0, 0, 0, // 9
				byte(OP_MULT),
				byte(OP_PRINT),
			},
		},
		{
			name: "simple true",
			in: []Statement{
				PrintStatment{
					Expr: LiteralExpression{TrueValue(true)},
				},
			},
			want: []byte{
				byte(OP_PUSH),
				byte(VAL_TRUE),
				byte(OP_PRINT),
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(tt *testing.T) {
			got, err := Compile(tst.in)
			if err != nil {
				tt.Fatalf("%s", err.Error())
			}
			if diff := cmp.Diff(tst.want, got); diff != "" {
				tt.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
