package lox

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]interface{}),
		enclosing: enclosing,
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) interface{} {
	if v, ok := e.values[name.lexeme]; ok {
		return v
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	panic(NewRuntimeError(name, "Undefined variable '"+name.lexeme+"'."))
}

func (e *Environment) getAt(distance int, name string) interface{} {
	return e.ancestor(distance).values[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) assign(name *Token, value interface{}) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}
	panic(NewRuntimeError(name, "Undefined variable '"+name.lexeme+"'."))
}

func (e *Environment) assignAt(distance int, name *Token, value interface{}) {
	e.ancestor(distance).values[name.lexeme] = value
}
