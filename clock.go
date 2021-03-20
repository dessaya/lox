package lox

import (
	"time"
)

type Clock struct{}

func (c Clock) Arity() int { return 0 }

func (c Clock) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().Unix())
}

func (c Clock) String() string { return "<native fn>" }
