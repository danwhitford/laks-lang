package laks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompile(t *testing.T) {
	var tests = []struct {
		name string
		in   []Expression
		want []byte
	}{
		{
			name: "literal",
			in: []Expression{
				{T: E_LIT, Value: int64(14)},
			},
			want: []byte{
				byte(OP_PUSH),
				14, 0, 0, 0, 0, 0, 0, 0, // 14
			},
		},
		{
			name: "expradd",
			in: []Expression{
				{
					T: E_BINOP,
					Value: BinaryExpression{
						BO_ADD,
						Expression{T: E_LIT, Value: int64(7)},
						Expression{T: E_LIT, Value: int64(9)},
					},
				},
			},
			want: []byte{
				byte(OP_PUSH),
				7, 0, 0, 0, 0, 0, 0, 0, // 14
				byte(OP_PUSH),
				9, 0, 0, 0, 0, 0, 0, 0, // 14
				byte(OP_ADD),
			},
		},
		{
			name: "print expression",
			in: []Expression{
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
			want: []byte{
				byte(OP_PUSH),
				7, 0, 0, 0, 0, 0, 0, 0, // 14
				byte(OP_PUSH),
				8, 0, 0, 0, 0, 0, 0, 0, // 14
				byte(OP_MULT),
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
