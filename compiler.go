package laks

import (
	"encoding/binary"
	"fmt"
)

type OpCode byte

const (
	OP_PUSH OpCode = iota
	OP_ADD
)

func Compile(exprs []Expression) ([]byte, error) {
	var bytecode []byte
	for _, expr := range exprs {
		b, err := compile_expr(expr)
		if err != nil {
			return bytecode, err
		}
		bytecode = append(bytecode, b...)
	}
	return bytecode, nil
}

func compile_expr(expr Expression) ([]byte, error) {
	switch expr.T {
	case E_LIT:
		return compile_lit(expr)
	case E_OP:
		return compile_binary_op(expr)
	default:
		return nil, fmt.Errorf("did not recognise expression type '%v'", expr.T)
	}
}

func compile_binary_op(expr Expression) ([]byte, error) {

	var bytecode []byte
	left, err := compile_expr(*expr.Left)
	if err != nil {
		return bytecode, err
	}
	bytecode = append(bytecode, left...)
	right, err := compile_expr(*expr.Right)
	if err != nil {
		return bytecode, err
	}
	bytecode = append(bytecode, right...)

	var op OpCode
	operator, ok := expr.Value.(BinaryOperator)
	if !ok {
		return nil, fmt.Errorf("expr '%v' did not have a valid BinaryOperator", expr)
	}
	switch operator {
	case BO_ADD:
		op = OP_ADD
	default:
		return nil, fmt.Errorf("dunno how to handle '%v'", operator)
	}

	return append(bytecode, byte(op)), nil
}

func compile_lit(expr Expression) ([]byte, error) {
	var buf []byte
	buf = append(buf, byte(OP_PUSH))
	val, ok := expr.Value.(int64)
	if !ok {
		return nil, fmt.Errorf("value of lit expr '%v' was not int64", expr)
	}
	return binary.Append(buf, binary.LittleEndian, val)
}
