package main

import (
	"fmt"
	"os"
	"path/filepath"

	core "github.com/salillakra/npp/core/interpreter"
	"github.com/salillakra/npp/frontend/lexer"
	"github.com/salillakra/npp/frontend/parser"
)

func main() {

	var filePath string
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	} else {
		fmt.Println("Please provide a file path as an argument.")
		return
	}

	fileExtension := filepath.Ext(filePath)

	if fileExtension != ".npp" {
		fmt.Println("Invalid file type. Please provide a .npp file.")
		return
	}

	dat, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(dat))
	p := parser.New(l, false) // Disabled debug output
	program := p.ParseProgram()
	i := core.New()
	i.Interpret(program)
}
