package main

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{make(map[string]interface{})}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) interface{} {
	if v, ok := e.values[name.lexeme]; ok {
		return v
	}
	panic(NewRuntimeError(name, "Undefined variable '"+name.lexeme+"'."))
}
func (e *Environment) assign(name *Token, value interface{}) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return
	}
	panic(NewRuntimeError(name, "Undefined variable '"+name.lexeme+"'."))
}
