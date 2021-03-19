package main

import "fmt"

type astPrinter struct {
	s string
}

func ExprToString(e Expr) string {
	a := &astPrinter{}
	e.accept(a)
	return a.s
}

func (a *astPrinter) visitBinary(b *Binary) {
	a.parenthesize(b.operator.lexeme, b.left, b.right)
}

func (a *astPrinter) visitGrouping(g *Grouping) {
	a.parenthesize("group", g.expression)
}

func (a *astPrinter) visitLiteral(l *Literal) {
	if l.value == nil {
		a.s = "nil"
		return
	}
	a.s = fmt.Sprintf("%v", l.value)
}

func (a *astPrinter) visitUnary(u *Unary) {
	a.parenthesize(u.operator.lexeme, u.right)
}

func (a *astPrinter) parenthesize(name string, exprs ...Expr) {
	a.s = "(" + name
	for _, expr := range exprs {
		a.s += " "
		aa := &astPrinter{}
		expr.accept(aa)
		a.s += aa.s
	}
	a.s += ")"
}
