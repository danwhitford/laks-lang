package laks

import (
	"encoding/binary"
	"fmt"
	"io"
)

type stack []int64

func (s *stack) push(i int64) {
	*s = append(*s, i)
}

func (s *stack) pop() int64 {
	v := (*s)[len((*s))-1]
	newlength := int(len(*s) - 1)
	*s = (*s)[:newlength]
	return v
}

type bytecode_interpreter struct {
	ip        int
	bytecode  []byte
	w         io.Writer
	val_stack stack
}

func Run(bytecode []byte, w io.Writer) error {
	bi := bytecode_interpreter{bytecode: bytecode, w: w}
	return bi.run()
}

func (bi *bytecode_interpreter) run() error {
	for bi.ip < len(bi.bytecode) {
		code_id := bi.read()
		switch code_id {
		case byte(OP_PUSH):
			bi.push_val()
		case byte(OP_MULT):
			bi.mult()
		case byte(OP_PRINT):
			bi.print()
		case byte(OP_ADD):
			bi.add()
		case byte(OP_DIV):
			bi.div()
		case byte(OP_MINUS):
			bi.minus()
		default:
			return fmt.Errorf("could not decode byte code '%v'", code_id)
		}
	}
	return nil
}

func (bi *bytecode_interpreter) minus() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()
	bi.val_stack.push(b - a)
}

func (bi *bytecode_interpreter) div() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()
	if a == 0 {
		panic("divide by zero")
	}
	bi.val_stack.push(b / a)
}

func (bi *bytecode_interpreter) add() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()
	bi.val_stack.push(a + b)}

func (bi *bytecode_interpreter) print() {
	fmt.Fprintf(bi.w, "%v\n", bi.val_stack.pop())
}

func (bi *bytecode_interpreter) mult() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()
	bi.val_stack.push(a * b)
}

func (bi *bytecode_interpreter) push_val() {
	var d int64
	read, err := binary.Decode(bi.bytecode[bi.ip:], binary.LittleEndian, &d)
	if err != nil {
		panic(err)
	}
	bi.ip += read
	bi.val_stack = append(bi.val_stack, d)
}

func (bi *bytecode_interpreter) read() byte {
	b := bi.bytecode[bi.ip]
	bi.ip++
	return b
}
