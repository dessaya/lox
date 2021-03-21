package lox

type returnSignal struct {
	value interface{}
}

type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

func NewLoxFunction(declaration *Function, closure *Environment) *LoxFunction {
	return &LoxFunction{declaration, closure}
}

func (f *LoxFunction) Arity() int { return len(f.declaration.params) }

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := NewEnvironment(f.closure)
	for i := 0; i < len(f.declaration.params); i++ {
		environment.define(f.declaration.params[i].lexeme, arguments[i])
	}
	var returnValue interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r, ok := r.(returnSignal); ok {
					returnValue = r.value
				} else {
					panic(r)
				}
			}
		}()
		interpreter.executeBlock(f.declaration.body, environment)
	}()
	return returnValue
}

func (f *LoxFunction) String() string { return "<function" + f.declaration.name.lexeme + ">" }
