package laks

import (
	// "fmt"
	"io"
)

func RunBytes(b []byte, w io.Writer) error {
	tokens, err := Tokenise(b)
	if err != nil {
		return err
	}

	// fmt.Printf("\t%v\n", tokens)

	exprs, err := Parse(tokens)
	if err != nil {
		return err
	}

	// for _, e := range exprs {
	// 	fmt.Printf("\t%v\n", e)
	// }

	bytecode, err := Compile(exprs)
	if err != nil {
		return err
	}

	// for _, b := range bytecode {
	// 	fmt.Printf("%x\n", b)
	// }
	// fmt.Println()

	err = Run(bytecode, w)
	if err != nil {
		return err
	}

	return nil
}
