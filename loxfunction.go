package lox

type LoxFunction struct {
	declaration *Function
}

func NewLoxFunction(declaration *Function) *LoxFunction {
	return &LoxFunction{declaration}
}

func (f *LoxFunction) Arity() int { return len(f.declaration.params) }

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := NewEnvironment(interpreter.globals)
	for i := 0; i < len(f.declaration.params); i++ {
		environment.define(f.declaration.params[i].lexeme, arguments[i])
	}
	interpreter.executeBlock(f.declaration.body, environment)
	return nil
}

func (f *LoxFunction) String() string { return "<function" + f.declaration.name.lexeme + ">" }
