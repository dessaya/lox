package lox

import "fmt"

type Token struct {
	kind    TokenType
	lexeme  string
	literal interface{}
	line    int
}

func NewToken(kind TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{kind, lexeme, literal, line}
}

func (t *Token) String() string {
	return fmt.Sprintf("%d %s %+v", t.kind, t.lexeme, t.literal)
}
