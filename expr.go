package main

type Expr interface {
	accept(v ExprVisitor) interface{}
}

type Binary struct {
	left     Expr
	operator *Token
	right    Expr
}

func NewBinary(left Expr, operator *Token, right Expr) *Binary {
	return &Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (b *Binary) accept(ev ExprVisitor) interface{} {
	return ev.visitBinary(b)
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{
		expression: expression,
	}
}

func (g *Grouping) accept(ev ExprVisitor) interface{} {
	return ev.visitGrouping(g)
}

type Literal struct {
	value interface{}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{
		value: value,
	}
}

func (l *Literal) accept(ev ExprVisitor) interface{} {
	return ev.visitLiteral(l)
}

type Unary struct {
	operator *Token
	right    Expr
}

func NewUnary(operator *Token, right Expr) *Unary {
	return &Unary{
		operator: operator,
		right:    right,
	}
}

func (u *Unary) accept(ev ExprVisitor) interface{} {
	return ev.visitUnary(u)
}

type ExprVisitor interface {
	visitBinary(b *Binary) interface{}
	visitGrouping(g *Grouping) interface{}
	visitLiteral(l *Literal) interface{}
	visitUnary(u *Unary) interface{}
}
