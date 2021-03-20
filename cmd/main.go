package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/dessaya/lox"
)

var interpreter = lox.NewInterpreter()

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
	if lox.HadError {
		os.Exit(65)
	}
	if lox.HadRuntimeError {
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
		lox.HadError = false
		if err == io.EOF {
			break
		}
	}
}

func run(source string) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := lox.NewParser(tokens)
	statements := parser.Parse()

	// Stop if there was a syntax error.
	if lox.HadError {
		return
	}

	interpreter.Interpret(statements)
}
