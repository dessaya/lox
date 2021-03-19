package main

import (
	"errors"
	"fmt"
)

type RuntimeError struct {
	error
	Token *Token
}

func NewRuntimeError(token *Token, msg string) RuntimeError {
	return RuntimeError{error: errors.New(msg), Token: token}
}

type Interpreter struct {
}

func (i *Interpreter) Interpret(expression Expr) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(RuntimeError); ok {
				ReportRuntimeError(e)
				return
			}
			panic(r)
		}
	}()

	value := i.evaluate(expression)
	fmt.Println(stringify(value))
}

func (i *Interpreter) visitBinary(b *Binary) interface{} {
	left := i.evaluate(b.left)
	right := i.evaluate(b.right)
	switch b.operator.kind {
	case GREATER:
		left, right := checkNumbers(b.operator, left, right)
		return left > right
	case GREATER_EQUAL:
		left, right := checkNumbers(b.operator, left, right)
		return left >= right
	case LESS:
		left, right := checkNumbers(b.operator, left, right)
		return left < right
	case LESS_EQUAL:
		left, right := checkNumbers(b.operator, left, right)
		return left <= right
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case MINUS:
		left, right := checkNumbers(b.operator, left, right)
		return left - right
	case PLUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right
			}
		}
		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right
			}
		}
		panic(NewRuntimeError(b.operator, "Operands must be two numbers or two strings."))
	case SLASH:
		left, right := checkNumbers(b.operator, left, right)
		return left / right
	case STAR:
		left, right := checkNumbers(b.operator, left, right)
		return left * right
	}
	panic("unreachable")
}

func (i *Interpreter) visitGrouping(g *Grouping) interface{} {
	return i.evaluate(g.expression)
}

func (i *Interpreter) visitLiteral(l *Literal) interface{} {
	return l.value
}

func (i *Interpreter) visitUnary(u *Unary) interface{} {
	right := i.evaluate(u.right)
	switch u.operator.kind {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		return -checkNumber(u.operator, right)
	}
	panic("unreachable")
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(i)
}

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	switch object := object.(type) {
	case bool:
		return object
	}
	return true
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func checkNumber(op *Token, x interface{}) float64 {
	if x, ok := x.(float64); ok {
		return x
	}
	panic(NewRuntimeError(op, "Operand must be a number."))
}

func checkNumbers(op *Token, left, right interface{}) (float64, float64) {
	if left, ok := left.(float64); ok {
		if right, ok := right.(float64); ok {
			return left, right
		}
	}
	panic(NewRuntimeError(op, "Operands must be numbers."))
}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	if s, ok := object.(string); ok {
		return fmt.Sprintf("%q", s)
	}
	return fmt.Sprintf("%v", object)
}
