package laks

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type=BinaryOperator
type BinaryOperator byte

const (
	BO_ADD BinaryOperator = iota
	BO_MINUS
	BO_MULT
	BO_DIV
	BO_EQ
)

type Statement any

type PrintStatment struct {
	Expr Statement
}

type BinaryExpression struct {
	Op    BinaryOperator
	Left  Statement
	Right Statement
}

type LiteralExpression struct {
	Value Value
}

func Parse(tokens []Token) ([]Statement, error) {
	p := parser{tokens: tokens}
	return p.parse()
}

type parser struct {
	tokens []Token
	curr   int
}

func (p *parser) parse() ([]Statement, error) {
	var exprs []Statement
	for p.curr < len(p.tokens) {
		expr, err := p.parse_statement()
		if err != nil {
			return exprs, err
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}

func (p *parser) parse_statement() (Statement, error) {
	t := p.peek()
	var stmt Statement
	var err error
	switch t.T {
	case T_INT:
		stmt, err = p.parse_bools()
	case T_KEYWORD:
		stmt, err = p.parse_keyword()
	default:
		return stmt, fmt.Errorf("do not know how to handle '%#v'", t)
	}

	if err != nil {
		return stmt, fmt.Errorf("error parsing statement. %s", err)
	}

	err = p.consume(T_SEMI)
	if err != nil {
		err = fmt.Errorf("error parsing statement. %v", err)
	}
	return stmt, err
}

func (p *parser) parse_keyword() (Statement, error) {
	kwd := p.read()
	switch kwd.Lexeme {
	case "print":
		expr, err := p.parse_bools()
		if err != nil {
			return nil, err
		}
		return PrintStatment{expr}, nil
	default:
		return nil, fmt.Errorf("do not recognise keyword '%v'", kwd.Lexeme)
	}
}

func (p *parser) parse_bools() (Statement, error) {
	expr, err := p.parse_expression()
	if err != nil {
		return expr, err
	}
	for p.peek().T == T_EQ_EQ {
		op_token := p.read()
		op := op_token_to_binary_op(op_token.T)
		r, err := p.parse_expression()
		if err != nil {
			return r, nil
		}
		expr = BinaryExpression{op, expr, r}
	}

	return expr, nil
}

func (p *parser) parse_expression() (Statement, error) {
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
		expr = BinaryExpression{op, expr, r}
	}

	return expr, nil
}

func (p *parser) parse_expression2() (Statement, error) {
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
		expr = BinaryExpression{op, expr, r}
	}

	return expr, nil
}

func (p *parser) parse_literal() (Statement, error) {
	t := p.read()
	switch t.T {
	case T_INT:
		d, err := strconv.ParseInt(t.Lexeme, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse literal '%s'. %s", t.Lexeme, err)
		}
		return LiteralExpression{IntValue(d)}, nil
	case T_KEYWORD:
		switch t.Lexeme {
		case "true":
			return LiteralExpression{TrueValue(true)}, nil
		case "false":
			return LiteralExpression{FalseValue(false)}, nil
		default:
			return nil, fmt.Errorf("could not parse literal as keyword '%#v'", t)
		}
	case T_STRING:
		return LiteralExpression{StringValue(t.Lexeme)}, nil
	default:
		return nil, fmt.Errorf("could not parse literal '%#v'", t)
	}
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
	case T_EQ_EQ:
		return BO_EQ
	default:
		panic("what is this '" + t.String() + "'")
	}
}

func (p *parser) consume(T TokenType) error {
	if p.curr >= len(p.tokens) {
		return fmt.Errorf("error consuming. EOF")
	}
	t := p.tokens[p.curr]
	if t.T != T {
		return fmt.Errorf("error consuming. wanted '%v' but got '%v'", T, t.T)
	}
	p.curr++
	return nil
}

func (p *parser) peek() Token {
	if p.curr >= len(p.tokens) {
		return Token{}
	}
	return p.tokens[p.curr]
}

func (p *parser) read() Token {
	t := p.tokens[p.curr]
	p.curr++
	return t
}
