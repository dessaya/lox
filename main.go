package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var hadError bool
var hadRuntimeError bool

var interpreter = &Interpreter{}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Printf("Usage: jlox [script]\n")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func loadFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func runFile(path string) {
	run(loadFile(path))
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		run(line)
		hadError = false
		if err == io.EOF {
			break
		}
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens)
	statements := parser.Parse()

	// Stop if there was a syntax error.
	if hadError {
		return
	}

	interpreter.Interpret(statements)
}

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
	hadRuntimeError = true
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}
