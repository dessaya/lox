package lox

import "fmt"

var (
	HadError        bool
	HadRuntimeError bool
)

func ReportError(line int, message string) {
	report(line, "", message)
}

func ReportTokenError(token *Token, message string) {
	if token.kind == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func ReportRuntimeError(err RuntimeError) {
	fmt.Printf("%s\n[line %d]\n", err.Error(), err.Token.line)
	HadRuntimeError = true
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	HadError = true
}
