package lox

// autogenerated with `python generate_ast.py .`

type Stmt interface {
	accept(v StmtVisitor) interface{}
}

type Block struct {
	statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		statements: statements,
	}
}

func (b *Block) accept(sv StmtVisitor) interface{} {
	return sv.visitBlockStmt(b)
}

type Expression struct {
	expression Expr
}

func NewExpression(expression Expr) *Expression {
	return &Expression{
		expression: expression,
	}
}

func (e *Expression) accept(sv StmtVisitor) interface{} {
	return sv.visitExpressionStmt(e)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) *If {
	return &If{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (i *If) accept(sv StmtVisitor) interface{} {
	return sv.visitIfStmt(i)
}

type Print struct {
	expression Expr
}

func NewPrint(expression Expr) *Print {
	return &Print{
		expression: expression,
	}
}

func (p *Print) accept(sv StmtVisitor) interface{} {
	return sv.visitPrintStmt(p)
}

type Var struct {
	name        *Token
	initializer Expr
}

func NewVar(name *Token, initializer Expr) *Var {
	return &Var{
		name:        name,
		initializer: initializer,
	}
}

func (v *Var) accept(sv StmtVisitor) interface{} {
	return sv.visitVarStmt(v)
}

type While struct {
	condition Expr
	body      Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
	return &While{
		condition: condition,
		body:      body,
	}
}

func (w *While) accept(sv StmtVisitor) interface{} {
	return sv.visitWhileStmt(w)
}

type StmtVisitor interface {
	visitBlockStmt(b *Block) interface{}
	visitExpressionStmt(e *Expression) interface{}
	visitIfStmt(i *If) interface{}
	visitPrintStmt(p *Print) interface{}
	visitVarStmt(v *Var) interface{}
	visitWhileStmt(w *While) interface{}
}
