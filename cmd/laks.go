package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/danwhitford/laks"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(">>> ")
	for scanner.Scan() { // Read each line
		input := scanner.Bytes() // Get the scanned text
		tokens, err := laks.Tokenise(input)
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

		fmt.Print(">>> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
