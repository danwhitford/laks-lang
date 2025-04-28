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
				Statement{
					T: ST_LIT,
					V: LiteralExpression{14},
				},
			},
			want: []byte{
				byte(OP_PUSH),
				14, 0, 0, 0, 0, 0, 0, 0, // 14
			},
		},
		{
			name: "expradd",
			in: []Statement{
				Statement{
					T: ST_BINEXPR,
					V: BinaryExpression{
						Op: BO_ADD,
						Left: Statement{
							T: ST_LIT,
							V: LiteralExpression{7},
						},
						Right: Statement{
							T: ST_LIT,
							V: LiteralExpression{9},
						},
					},
				},
			},
			want: []byte{
				byte(OP_PUSH),
				7, 0, 0, 0, 0, 0, 0, 0, // 7
				byte(OP_PUSH),
				9, 0, 0, 0, 0, 0, 0, 0, // 9
				byte(OP_ADD),
			},
		},
		{
			name: "print expression",
			in: []Statement{
				Statement{
					T: ST_PRINT,
					V: PrintStatment{
						Expr: Statement{
							T: ST_BINEXPR,
							V: BinaryExpression{
								Op: BO_MULT,
								Left: Statement{
									T: ST_LIT,
									V: LiteralExpression{7},
								},
								Right: Statement{
									T: ST_LIT,
									V: LiteralExpression{9},
								},
							},
						},
					},
				},
			},
			want: []byte{
				byte(OP_PUSH),
				7, 0, 0, 0, 0, 0, 0, 0, // 7
				byte(OP_PUSH),
				9, 0, 0, 0, 0, 0, 0, 0, // 9
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
