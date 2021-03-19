// +build ignore

package main

import "fmt"

type astPrinter struct {
}

func ExprToString(e Expr) string {
	a := astPrinter{}
	return e.accept(a).(string)
}

func (a astPrinter) visitBinaryExpr(b *Binary) interface{} {
	return a.parenthesize(b.operator.lexeme, b.left, b.right)
}

func (a astPrinter) visitGroupingExpr(g *Grouping) interface{} {
	return a.parenthesize("group", g.expression)
}

func (a astPrinter) visitLiteralExpr(l *Literal) interface{} {
	if l.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", l.value)
}

func (a astPrinter) visitUnaryExpr(u *Unary) interface{} {
	return a.parenthesize(u.operator.lexeme, u.right)
}

func (a astPrinter) parenthesize(name string, exprs ...Expr) string {
	s := "(" + name
	for _, expr := range exprs {
		s += " "
		aa := astPrinter{}
		ss := expr.accept(aa).(string)
		s += ss
	}
	s += ")"
	return s
}
