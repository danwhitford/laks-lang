package laks

import (
	"encoding/binary"
	"fmt"
)

//go:generate stringer -type=OpCode
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
	case E_BINOP:
		return compile_binary_op(expr)
	default:
		return nil, fmt.Errorf("did not recognise expression type '%v'", expr.T)
	}
}

func compile_binary_op(expr Expression) ([]byte, error) {
	var bytecode []byte

	bexpr, ok := expr.Value.(BinaryExpression)
	if !ok {
		return bytecode, fmt.Errorf("failed to convert '%v' to a BinaryExpression", expr)
	}
	left, err := compile_expr(bexpr.Left)
	if err != nil {
		return bytecode, err
	}
	bytecode = append(bytecode, left...)
	right, err := compile_expr(bexpr.Right)
	if err != nil {
		return bytecode, err
	}
	bytecode = append(bytecode, right...)

	var op OpCode
	switch bexpr.Operator {
	case BO_ADD:
		op = OP_ADD
	default:
		return nil, fmt.Errorf("dunno how to handle '%v'", bexpr.Operator)
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
