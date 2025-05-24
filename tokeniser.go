package laks

import (
	"fmt"
	"slices"
	"strings"
)

//go:generate stringer -type=TokenType
type TokenType byte

const (
	T_INT TokenType = iota
	T_SEMI
	T_MULT
	T_ADD
	T_DIV
	T_MINUS
	T_KEYWORD
	T_EQ
	T_EQ_EQ
	T_STRING
)

type Token struct {
	T      TokenType
	Lexeme string
}

type tokeniser struct {
	src     []byte
	current int
	tokens  []Token
}

func Tokenise(src []byte) ([]Token, error) {
	var t = tokeniser{src: src}
	err := t.tokenise()
	return t.tokens, err
}

func (t *tokeniser) tokenise() error {
	for t.current < len(t.src) {
		r := t.peek()

		if r < '!' {
			t.read()
			continue
		}

		if r >= '0' && r <= '9' {
			t.tokenise_number()
		} else if slices.Contains([]byte{'*', '+', '/', '-', '='}, r) {
			t.tokenise_operator()
		} else if r == ';' {
			t.read()
			t.tokens = append(t.tokens, Token{T_SEMI, string(r)})
		} else if r >= 'a' && r <= 'z' {
			t.tokenise_keyword()
		} else if r == '#' {
			t.eat_comment()
		} else if r == '"' {
			t.tokenise_string()
		} else {
			return fmt.Errorf("cannot tokenise '%c'", r)
		}
	}

	return nil
}

func (t *tokeniser) tokenise_string() {
	t.read() // The opening quotes

	var sb strings.Builder

	for t.current < len(t.src) {
		r := t.read()
		if r == '"' {
			t.tokens = append(t.tokens, Token{T_STRING, sb.String()})
			return
		}
		sb.WriteByte(r)
	}

	panic("got to end of file while reading string")
}

func (t *tokeniser) eat_comment() {
	for t.current < len(t.src) {
		r := t.peek()

		if r == '\n' {
			break
		}

		t.read()
	}
}

func (t *tokeniser) tokenise_keyword() {
	var sb strings.Builder

	for t.current < len(t.src) {
		r := t.peek()

		if r >= 'a' && r <= 'z' {
			sb.WriteByte(t.read())
		} else {
			break
		}
	}

	t.tokens = append(t.tokens, Token{T_KEYWORD, sb.String()})
}

func (t *tokeniser) tokenise_operator() {
	r := t.read()
	switch r {
	case '*':
		t.tokens = append(t.tokens, Token{T_MULT, string(r)})
	case '+':
		t.tokens = append(t.tokens, Token{T_ADD, string(r)})
	case '-':
		t.tokens = append(t.tokens, Token{T_MINUS, string(r)})
	case '/':
		t.tokens = append(t.tokens, Token{T_DIV, string(r)})
	case '=':
		if t.peek() == '=' {
			t.read()
			t.tokens = append(t.tokens, Token{T_EQ_EQ, "=="})
		} else {
			t.tokens = append(t.tokens, Token{T_EQ, string(r)})
		}
	}
}

func (t *tokeniser) peek() byte {
	return t.src[t.current]
}

func (t *tokeniser) tokenise_number() {
	var sb strings.Builder

	for t.current < len(t.src) {
		r := t.peek()

		if r >= '0' && r <= '9' {
			sb.WriteByte(t.read())
		} else {
			break
		}
	}

	t.tokens = append(t.tokens, Token{T_INT, sb.String()})
}

func (t *tokeniser) read() byte {
	r := t.src[t.current]
	t.current++
	return r
}
