package main

import (
	"os"
	"testing"

	core "github.com/salillakra/npp/core/interpreter"
	"github.com/salillakra/npp/frontend/lexer"
	"github.com/salillakra/npp/frontend/parser"
)

func TestCodeNPP(t *testing.T) {
	code, err := os.ReadFile("./hello.npp")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if len(code) == 0 {
		t.Fatal("File is empty")
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	l := lexer.New(string(code))
	p := parser.New(l, false) // Disabled debug output
	program := p.ParseProgram()
	i := core.New()
	i.Interpret(program)

	w.Close()
	os.Stdout = oldStdout
	var buf [4096]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])

	expected := "2\nhello world\nfuck off!\n"
	if output != expected {
		t.Errorf("Unexpected output.\nGot:\n%q\nWant:%q", output, expected)
	}
}
