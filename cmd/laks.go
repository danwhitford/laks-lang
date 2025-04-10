package main

import (
	"fmt"
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

	tokens, err := laks.Tokenise(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\t%v\n", tokens)
	exprs, err := laks.Parse(tokens)
	if err != nil {
		panic(err)
	}
	for _, e := range exprs {
		fmt.Printf("\t%v\n", e)
	}

	bytecode, err := laks.Compile(exprs)
	if err != nil {
		panic(err)
	}

	err = laks.Run(bytecode, os.Stdout)
	if err != nil {
		panic(err)
	}
}
