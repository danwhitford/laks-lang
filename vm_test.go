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
		// {
		// 	name: "print4",
		// 	in: 
		// }

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
