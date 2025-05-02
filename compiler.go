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
	OP_DIV
	OP_MINUS
	OP_EQ
)

func compileLiteralExpression(expr LiteralExpression) ([]byte, error) {
	var buf []byte
	var err error
	switch expr.Value.T {
	case VAL_INT:
		buf = append(buf, byte(VAL_INT))
		buf = append(buf, byte(OP_PUSH))
		buf, err = binary.Append(buf, binary.LittleEndian, expr.Value.Val.(int64))
		if err != nil {
			err = fmt.Errorf("error appending '%#v'. %v", expr.Value, err)
		}
	case VAL_TRUE:
		buf = append(buf, byte(OP_PUSH))
		buf = append(buf, byte(VAL_TRUE))
	case VAL_FALSE:
		buf = append(buf, byte(OP_PUSH))
		buf = append(buf, byte(VAL_FALSE))
	default:
		return buf, fmt.Errorf("do not know how to compile litexpr '%v'", expr)
	}
	return buf, err
}

func compileBinaryExpression(bexpr BinaryExpression) ([]byte, error) {
	var buf []byte
	left, err := compileStatement(bexpr.Left)
	if err != nil {
		return buf, err
	}
	buf = append(buf, left...)
	right, err := compileStatement(bexpr.Right)
	if err != nil {
		return buf, err
	}
	buf = append(buf, right...)

	switch bexpr.Op {
	case BO_ADD:
		buf = append(buf, byte(OP_ADD))
	case BO_MULT:
		buf = append(buf, byte(OP_MULT))
	case BO_DIV:
		buf = append(buf, byte(OP_DIV))
	case BO_MINUS:
		buf = append(buf, byte(OP_MINUS))
	case BO_EQ:
		buf = append(buf, byte(OP_EQ))
	default:
		return buf, fmt.Errorf("unknown operator '%v'", bexpr.Op)
	}

	return buf, nil
}

func compilePrint(p PrintStatment) ([]byte, error) {
	b, err := compileStatement(p.Expr)
	if err != nil {
		return b, fmt.Errorf("error compiling expression for printing '%v'. '%v'", p.Expr, err)
	}
	b = append(b, byte(OP_PRINT))
	return b, nil
}

func compileStatement(stmt Statement) ([]byte, error) {
	switch stmt.T {
	case ST_PRINT:
		return compilePrint(stmt.V.(PrintStatment))
	case ST_BINEXPR:
		return compileBinaryExpression(stmt.V.(BinaryExpression))
	case ST_LIT:
		return compileLiteralExpression(stmt.V.(LiteralExpression))
	default:
		return nil, fmt.Errorf("unknown statement type '%v'", stmt.T)
	}
}

func Compile(stmts []Statement) ([]byte, error) {
	var bytecode []byte
	for _, stmt := range stmts {
		b, err := compileStatement(stmt)
		if err != nil {
			return bytecode, fmt.Errorf("error compiling statement '%v'. '%v'", stmt, err)
		}
		bytecode = append(bytecode, b...)
	}
	return bytecode, nil
}
