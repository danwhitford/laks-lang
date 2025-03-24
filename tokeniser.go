package laks

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
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
)

type Token struct {
	T      TokenType
	Lexeme string
}

type tokeniser struct {
	runes   []rune
	current int
	tokens  []Token
}

func Tokenise(src string) ([]Token, error) {
	var runes = []rune(src)
	var t = tokeniser{runes: runes}
	err := t.tokenise()
	return t.tokens, err
}

func (t *tokeniser) tokenise() error {
	for t.current < len(t.runes) {
		r := t.peek()

		if unicode.IsSpace(r) {
			t.read()
			continue
		}

		if unicode.IsDigit(r) {
			t.tokenise_number()
		} else if slices.Contains([]rune{'*', '+', '/', '-'}, r) {
			t.tokenise_operator()
		} else if r == ';' {
			t.read()
			t.tokens = append(t.tokens, Token{T_SEMI, string(r)})
		} else {
			return fmt.Errorf("cannot tokenise '%c'", r)
		}
	}

	return nil
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
	}
}

func (t *tokeniser) peek() rune {
	return t.runes[t.current]
}

func (t *tokeniser) tokenise_number() {
	var sb strings.Builder

	for t.current < len(t.runes) {
		r := t.peek()

		if unicode.IsDigit(r) {
			sb.WriteRune(t.read())
		} else {
			break
		}
	}

	t.tokens = append(t.tokens, Token{T_INT, sb.String()})
}

func (t *tokeniser) read() rune {
	r := t.runes[t.current]
	t.current++
	return r
}
