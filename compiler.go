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
	OP_MULT
	OP_PRINT
)

func (expr LiteralExpression) Compile() ([]byte, error) {
	var buf []byte
	buf = append(buf, byte(OP_PUSH))
	return binary.Append(buf, binary.LittleEndian, expr.Value)
}

func (expr BinaryExpression) Compile() ([]byte, error) {
	var buf []byte
	left, err := expr.Left.Compile()
	if err != nil {
		return buf, err
	}
	buf = append(buf, left...)
	right, err := expr.Right.Compile()
	if err != nil {
		return buf, err
	}
	buf = append(buf, right...)
	
	switch expr.Op {
	case BO_ADD:
		buf = append(buf, byte(OP_ADD))
	case BO_MULT:
		buf = append(buf, byte(OP_MULT))
	default:
		return buf, fmt.Errorf("unknown operator '%v'", expr.Op)
	}

	return buf, nil	
}

func (expr PrintStatment) Compile() ([]byte, error) {
	b, err := expr.Expr.Compile()
	if err != nil {
		return b, fmt.Errorf("error compiling expression for printing '%v'. '%v'", expr, err)
	}	
	b = append(b, byte(OP_PRINT))
	return b, nil
}

func Compile(stmts []Statement) ([]byte, error) {
	var bytecode []byte
	for _, stmt := range stmts {
		b, err := stmt.Compile()
		if err != nil {
			return bytecode, err
		}
		bytecode = append(bytecode, b...)
	}
	return bytecode, nil
}
