package lox

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
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	globals.define("clock", Clock{})
	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[Expr]int),
	}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(RuntimeError); ok {
				ReportRuntimeError(e)
				return
			}
			panic(r)
		}
	}()

	for _, s := range statements {
		i.execute(s)
	}
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.accept(i)
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment

	i.environment = environment
	defer func() { i.environment = previous }()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) visitBlockStmt(stmt *Block) interface{} {
	i.executeBlock(stmt.statements, NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) visitExpressionStmt(stmt *Expression) interface{} {
	i.evaluate(stmt.expression)
	return nil
}

func (i *Interpreter) visitFunctionStmt(stmt *Function) interface{} {
	function := NewLoxFunction(stmt, i.environment)
	i.environment.define(stmt.name.lexeme, function)
	return nil
}

func (i *Interpreter) visitIfStmt(stmt *If) interface{} {
	if isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
	return nil
}

func (i *Interpreter) visitPrintStmt(stmt *Print) interface{} {
	value := i.evaluate(stmt.expression)
	fmt.Printf("%s\n", stringify(value))
	return nil
}

func (i *Interpreter) visitReturnStmt(stmt *Return) interface{} {
	var value interface{}
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}
	panic(returnSignal{value})
}

func (i *Interpreter) visitVarStmt(stmt *Var) interface{} {
	var value interface{}
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}
	i.environment.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) visitWhileStmt(stmt *While) interface{} {
	for isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}
	return nil
}

func (i *Interpreter) visitAssignExpr(expr *Assign) interface{} {
	value := i.evaluate(expr.value)

	if distance, ok := i.locals[expr]; ok {
		i.environment.assignAt(distance, expr.name, value)
	} else {
		i.globals.assign(expr.name, value)
	}

	return value
}

func (i *Interpreter) visitBinaryExpr(b *Binary) interface{} {
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

func (i *Interpreter) visitCallExpr(expr *Call) interface{} {
	callee := i.evaluate(expr.callee)

	var arguments []interface{}
	for _, argument := range expr.arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	if function, ok := callee.(LoxCallable); ok {
		if len(arguments) != function.Arity() {
			panic(NewRuntimeError(expr.paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments))))
		}
		return function.Call(i, arguments)
	}
	panic(NewRuntimeError(expr.paren, "Can only call functions and classes."))
}

func (i *Interpreter) visitGroupingExpr(g *Grouping) interface{} {
	return i.evaluate(g.expression)
}

func (i *Interpreter) visitLiteralExpr(l *Literal) interface{} {
	return l.value
}

func (i *Interpreter) visitLogicalExpr(expr *Logical) interface{} {
	left := i.evaluate(expr.left)

	if expr.operator.kind == OR {
		if isTruthy(left) {
			return left
		}
	} else { // AND
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.right)
}

func (i *Interpreter) visitUnaryExpr(u *Unary) interface{} {
	right := i.evaluate(u.right)
	switch u.operator.kind {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		return -checkNumber(u.operator, right)
	}
	panic("unreachable")
}

func (i *Interpreter) visitVariableExpr(expr *Variable) interface{} {
	return i.lookUpVariable(expr.name, expr)
}

func (i *Interpreter) lookUpVariable(name *Token, expr Expr) interface{} {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.getAt(distance, name.lexeme)
	}
	return i.globals.get(name)
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
		return s
	}
	return fmt.Sprintf("%v", object)
}
