package main

import (
	"io"
	"os"

	"github.com/danwhitford/laks"
)

func main() {
	var f *os.File

	if len(os.Args) > 1 {
		fname := os.Args[1]
		var err error
		f, err = os.Open(fname)
		if err != nil {
			panic(err)
		}
	} else {
		f = os.Stdin
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	laks.RunBytes(b, os.Stdout)
}
