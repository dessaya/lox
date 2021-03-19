package main

import "errors"

type ParseError error

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(RuntimeError); ok {
				p.synchronize()
				return
			}
			panic(r)
		}
	}()

	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return NewVar(name, initializer)
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return NewPrint(value)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return NewExpression(expr)
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) match(types ...TokenType) bool {
	for _, kind := range types {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(kind TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().kind == kind
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().kind == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return NewLiteral(false)
	case p.match(TRUE):
		return NewLiteral(true)
	case p.match(NIL):
		return NewLiteral(nil)
	case p.match(NUMBER, STRING):
		return NewLiteral(p.previous().literal)
	case (p.match(IDENTIFIER)):
		return NewVariable(p.previous())
	case p.match(LEFT_PAREN):
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	}
	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) consume(kind TokenType, message string) *Token {
	if p.check(kind) {
		return p.advance()
	}
	panic(p.error(p.peek(), message))
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().kind == SEMICOLON {
			return
		}

		switch p.peek().kind {
		case CLASS:
			return
		case FUN:
			return
		case VAR:
			return
		case FOR:
			return
		case IF:
			return
		case WHILE:
			return
		case PRINT:
			return
		case RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) error(token *Token, message string) ParseError {
	ReportTokenError(token, message)
	return ParseError(errors.New(message))
}
