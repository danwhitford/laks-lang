package laks

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type=ExpressionType
type ExpressionType byte

const (
	E_LIT ExpressionType = iota
	E_OP
)

//go:generate stringer -type=BinaryOperator
type BinaryOperator byte

const (
	BO_ADD BinaryOperator = iota
	BO_MINUS
	BO_MULT
	BO_DIV
)

type Expression struct {
	T     ExpressionType
	Value any
	Left  *Expression
	Right *Expression
}

func Parse(tokens []Token) ([]Expression, error) {
	p := parser{tokens: tokens}
	return p.parse()
}

type parser struct {
	tokens []Token
	curr   int
}

func (p *parser) parse() ([]Expression, error) {
	var exprs []Expression
	for p.curr < len(p.tokens) {
		expr, err := p.parse_statement()
		if err != nil {
			return exprs, err
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}

func (p *parser) parse_statement() (Expression, error) {
	t := p.peek()
	var stmt Expression
	var err error
	switch t.T {
	case T_INT:
		stmt, err = p.parse_expression()
	default:
		return Expression{}, fmt.Errorf("do not know how to handle '%#v'", t)
	}

	if err != nil {
		return Expression{}, fmt.Errorf("error parsing statement. %s", err)
	}

	err = p.consume(T_SEMI)
	return stmt, err
}

func (p *parser) parse_expression() (Expression, error) {
	expr, err := p.parse_expression2()
	if err != nil {
		return expr, err
	}
	for p.peek().T == T_ADD || p.peek().T == T_MINUS {
		op_token := p.read()
		op := op_token_to_binary_op(op_token.T)
		r, err := p.parse_expression2()
		if err != nil {
			return r, nil
		}
		expr = Expression{
			E_OP,
			op,
			&Expression{expr.T, expr.Value, expr.Left, expr.Right},
			&r,
		}
	}

	return expr, nil
}

func (p *parser) parse_expression2() (Expression, error) {
	expr, err := p.parse_literal()
	if err != nil {
		return expr, err
	}
	for p.peek().T == T_MULT || p.peek().T == T_DIV {
		op_token := p.read()
		op := op_token_to_binary_op(op_token.T)
		r, err := p.parse_literal()
		if err != nil {
			return r, nil
		}
		expr = Expression{
			E_OP,
			op,
			&Expression{expr.T, expr.Value, expr.Left, expr.Right},
			&r,
		}
	}

	return expr, nil
}

func op_token_to_binary_op(t TokenType) BinaryOperator {
	switch t {
	case T_ADD:
		return BO_ADD
	case T_DIV:
		return BO_DIV
	case T_MINUS:
		return BO_MINUS
	case T_MULT:
		return BO_MULT
	default:
		panic("what is this '" + t.String() + "'")
	}
}

func (p *parser) parse_literal() (Expression, error) {
	t := p.read()
	d, err := strconv.ParseInt(t.Lexeme, 10, 64)
	if err != nil {
		return Expression{}, fmt.Errorf("could not parse literal '%s'. %s", t.Lexeme, err)
	}
	return Expression{T: E_LIT, Value: d}, nil
}

func (p *parser) consume(T TokenType) error {
	t := p.tokens[p.curr]
	if t.T != T {
		return fmt.Errorf("error consuming. wanted '%v' but got '%v'", T, t.T)
	}
	p.curr++
	return nil
}

func (p *parser) peek() Token {
	return p.tokens[p.curr]
}

func (p *parser) read() Token {
	t := p.tokens[p.curr]
	p.curr++
	return t
}
