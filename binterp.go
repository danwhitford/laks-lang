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
	VAL_STRING
)

type Value any
type IntValue int64
type TrueValue bool
type FalseValue bool
type StringValue string

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

func (bi *bytecode_interpreter) minus() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()
	bi.val_stack.push(IntValue(b.(IntValue) - a.(IntValue)))
}

func (bi *bytecode_interpreter) div() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()
	if a.(IntValue) == 0 {
		panic("divide by zero")
	}
	bi.val_stack.push(IntValue(b.(IntValue) / a.(IntValue)))
}

func (bi *bytecode_interpreter) add() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()

	switch a.(type) {
	case IntValue:
		bi.val_stack.push(IntValue(a.(IntValue) + b.(IntValue)))
	case StringValue:
		bi.val_stack.push(StringValue(b.(StringValue) + a.(StringValue)))
	}

}

func (bi *bytecode_interpreter) print() {
	v := bi.val_stack.pop()
	switch v.(type) {
	case IntValue, StringValue:
		fmt.Fprintf(bi.w, "%v\n", v)
	case TrueValue:
		fmt.Fprintln(bi.w, "true")
	case FalseValue:
		fmt.Fprintln(bi.w, "false")
	default:
		fmt.Fprintf(bi.w, "%v\n", v)

	}
}

func (bi *bytecode_interpreter) mult() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()
	bi.val_stack.push(IntValue(a.(IntValue) * b.(IntValue)))
}

func (bi *bytecode_interpreter) eq() {
	a := bi.val_stack.pop()
	b := bi.val_stack.pop()

	if a == b {
		bi.val_stack.push(TrueValue(true))
	} else {
		bi.val_stack.push(FalseValue(false))
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
		bi.val_stack = append(bi.val_stack, IntValue(int64(d)))
	case byte(VAL_TRUE):
		bi.val_stack = append(bi.val_stack, TrueValue(true))
	case byte(VAL_FALSE):
		bi.val_stack = append(bi.val_stack, FalseValue(false))
	case byte(VAL_STRING):
		start := bi.ip
		end := start
		for {
			b := bi.read()
			if b == 0 {
				break
			}
			end++
		}
		bi.val_stack = append(bi.val_stack, StringValue(string(bi.bytecode[start:end])))
	default:
		panic(fmt.Sprintf("Could not convert '%v' to ValueType", val_byte))
	}
}

func (bi *bytecode_interpreter) read() byte {
	b := bi.bytecode[bi.ip]
	bi.ip++
	return b
}
