package laks

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_VM(t *testing.T) {
	var tests = []struct {
		name string
		in   []byte
		want string
	}{
		{
			name: "print a product",
			in: []byte{
				byte(OP_PUSH),
				7, 0, 0, 0, 0, 0, 0, 0, // 14
				byte(OP_PUSH),
				8, 0, 0, 0, 0, 0, 0, 0, // 14
				byte(OP_MULT),
				byte(OP_PRINT),
			},
			want: "56\n",
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(tt *testing.T) {
			var w bytes.Buffer
			err := Run(tst.in, &w)
			if err != nil {
				tt.Fatalf("%s", err.Error())
			}
			got := w.String()

			if diff := cmp.Diff(tst.want, got); diff != "" {
				tt.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStac(t *testing.T) {
	var s stack
	var i int64
	for i = range 5 {
		s.push(i)
	}

	if s.pop() != 4 {
		t.Errorf("ohno")
	}
	if s.pop() != 3 {
		t.Errorf("ohno")
	}
	if s.pop() != 2 {
		t.Errorf("ohno")
	}
	if s.pop() != 1 {
		t.Errorf("ohno")
	}
	if s.pop() != 0 {
		t.Errorf("ohno")
	}
}
