package laks

import (
	"encoding/binary"
	"fmt"
	"io"
)

//go:generate stringer -type=ValueType
type ValueType byte

const (
	VAL_INT ValueType = iota
	VAL_TRUE
	VAL_FALSE
)

type Value struct {
	T   ValueType
	Val any
}

type stack []Value

func (s *stack) push(i Value) {
	*s = append(*s, i)
}

func (s *stack) pop() Value {
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
		case byte(OP_EQ):
			bi.eq()
		default:
			return fmt.Errorf("could not decode byte code '%v'", code_id)
		}
	}
	return nil
}

func checkType(want, got ValueType) {
	if want != got {
		panic(fmt.Sprintf("Wanted '%v' but got '%v'", want, got))
	}
}

func (bi *bytecode_interpreter) minus() {
	a := bi.val_stack.pop()
	checkType(VAL_INT, a.T)
	b := bi.val_stack.pop()
	checkType(VAL_INT, b.T)
	bi.val_stack.push(Value{VAL_INT, b.Val.(int64) - a.Val.(int64)})
}

func (bi *bytecode_interpreter) div() {
	a := bi.val_stack.pop()
	checkType(VAL_INT, a.T)
	b := bi.val_stack.pop()
	checkType(VAL_INT, b.T)
	if a.Val.(int64) == 0 {
		panic("divide by zero")
	}
	bi.val_stack.push(Value{VAL_INT, b.Val.(int64) / a.Val.(int64)})
}

func (bi *bytecode_interpreter) add() {
	a := bi.val_stack.pop()
	checkType(VAL_INT, a.T)
	b := bi.val_stack.pop()
	checkType(VAL_INT, b.T)
	bi.val_stack.push(Value{VAL_INT, a.Val.(int64) + b.Val.(int64)})
}

func (bi *bytecode_interpreter) print() {
	v := bi.val_stack.pop()
	switch v.T {
	case VAL_INT:
		fmt.Fprintf(bi.w, "%v\n", v.Val)
	case VAL_TRUE:
		fmt.Fprintln(bi.w, "true")
	case VAL_FALSE:
		fmt.Fprintln(bi.w, "false")
	default:
		fmt.Fprintf(bi.w, "%v\n", v)

	}
}

func (bi *bytecode_interpreter) mult() {
	a := bi.val_stack.pop()
	checkType(VAL_INT, a.T)
	b := bi.val_stack.pop()
	checkType(VAL_INT, b.T)
	bi.val_stack.push(Value{VAL_INT, a.Val.(int64) * b.Val.(int64)})
}

func (bi *bytecode_interpreter) eq() {
	a := bi.val_stack.pop()
	checkType(VAL_INT, a.T)
	b := bi.val_stack.pop()
	checkType(VAL_INT, b.T)
	if a == b {
		bi.val_stack.push(Value{VAL_TRUE, true})
	} else {
		bi.val_stack.push(Value{VAL_FALSE, false})
	}
}

func (bi *bytecode_interpreter) push_val() {
	val_byte := bi.read()
	switch val_byte {
	case byte(VAL_INT):
		var d int64
		read, err := binary.Decode(bi.bytecode[bi.ip:], binary.LittleEndian, &d)
		if err != nil {
			panic(err)
		}
		bi.ip += read
		bi.val_stack = append(bi.val_stack, Value{VAL_INT, d})
	case byte(VAL_TRUE):
		bi.val_stack = append(bi.val_stack, Value{VAL_TRUE, true})
	case byte(VAL_FALSE):
		bi.val_stack = append(bi.val_stack, Value{VAL_FALSE, false})
	default:
		panic(fmt.Sprintf("Could not convert '%v' to ValueType", val_byte))
	}

}

func (bi *bytecode_interpreter) read() byte {
	b := bi.bytecode[bi.ip]
	bi.ip++
	return b
}
