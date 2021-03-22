package lox

type FunctionType int

const (
	NONE = FunctionType(iota)
	FUNCTION
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

func NewResolver(i *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     i,
		scopes:          nil,
		currentFunction: NONE,
	}
}

func (r *Resolver) visitAssignExpr(a *Assign) interface{} {
	r.resolveExpr(a.value)
	r.resolveLocal(a, a.name)
	return nil
}

func (r *Resolver) visitBinaryExpr(b *Binary) interface{} {
	r.resolveExpr(b.left)
	r.resolveExpr(b.right)
	return nil
}

func (r *Resolver) visitCallExpr(c *Call) interface{} {
	r.resolveExpr(c.callee)
	for _, argument := range c.arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) visitGroupingExpr(g *Grouping) interface{} {
	r.resolveExpr(g.expression)
	return nil
}

func (r *Resolver) visitLiteralExpr(l *Literal) interface{} {
	return nil
}

func (r *Resolver) visitLogicalExpr(l *Logical) interface{} {
	r.resolveExpr(l.left)
	r.resolveExpr(l.right)
	return nil
}

func (r *Resolver) visitUnaryExpr(u *Unary) interface{} {
	r.resolveExpr(u.right)
	return nil
}

func (r *Resolver) visitVariableExpr(v *Variable) interface{} {
	if len(r.scopes) > 0 {
		if defined, ok := r.scopes[len(r.scopes)-1][v.name.lexeme]; ok && !defined {
			ReportTokenError(v.name, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(v, v.name)
	return nil
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) visitBlockStmt(b *Block) interface{} {
	r.beginScope()
	r.resolveStmts(b.statements)
	r.endScope()
	return nil
}

func (r *Resolver) Resolve(statements []Stmt) {
	r.resolveStmts(statements)
}

func (r *Resolver) resolveStmts(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.accept(r)
}

func (r *Resolver) resolveExpr(e Expr) {
	e.accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) visitExpressionStmt(e *Expression) interface{} {
	r.resolveExpr(e.expression)
	return nil
}

func (r *Resolver) visitFunctionStmt(f *Function) interface{} {
	r.declare(f.name)
	r.define(f.name)
	r.resolveFunction(f, FUNCTION)
	return nil
}

func (r *Resolver) resolveFunction(function *Function, kind FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = kind

	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(function.body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

func (r *Resolver) visitIfStmt(i *If) interface{} {
	r.resolveExpr(i.condition)
	r.resolveStmt(i.thenBranch)
	if i.elseBranch != nil {
		r.resolveStmt(i.elseBranch)
	}
	return nil
}

func (r *Resolver) visitPrintStmt(p *Print) interface{} {
	r.resolveExpr(p.expression)
	return nil
}

func (r *Resolver) visitReturnStmt(ret *Return) interface{} {
	if r.currentFunction == NONE {
		ReportTokenError(ret.keyword, "Can't return from top-level code.")
	}

	if ret.value != nil {
		r.resolveExpr(ret.value)
	}
	return nil
}

func (r *Resolver) declare(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.lexeme]; ok {
		ReportTokenError(name, "Already variable with this name in this scope.")
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = true
}

func (r *Resolver) visitVarStmt(v *Var) interface{} {
	r.declare(v.name)
	if v.initializer != nil {
		r.resolveExpr(v.initializer)
	}
	r.define(v.name)
	return nil
}

func (r *Resolver) visitWhileStmt(w *While) interface{} {
	r.resolveExpr(w.condition)
	r.resolveStmt(w.body)
	return nil
}
