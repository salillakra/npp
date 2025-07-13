package interpreter

import (
	"fmt"

	"github.com/salillakra/npp/frontend/lexer"
	"github.com/salillakra/npp/frontend/parser"
)

// Object represents a value in the language (number or string).
type Object interface {
	String() string
}

// IntObject represents an integer value.
type IntObject struct {
	Value int64
}

func (i *IntObject) String() string { return fmt.Sprintf("%d", i.Value) }

// StringObject represents a string value.
type StringObject struct {
	Value string
}

func (s *StringObject) String() string { return s.Value }

// Environment stores variable bindings.
type Environment struct {
	store map[string]Object
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

// Interpreter evaluates the AST.
type Interpreter struct {
	env *Environment
}

// New creates a new Interpreter with optional sassy comments.
func New() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(),
	}
}

// Interpret executes the program.
func (i *Interpreter) Interpret(program *parser.Program) {
	if program == nil || program.Statements == nil {
		return
	}
	for _, stmt := range program.Statements {
		if stmt != nil {
			i.evalStatement(stmt)
		}
	}
}

// evalStatement evaluates a statement.
func (i *Interpreter) evalStatement(stmt parser.Statement) {
	if stmt == nil {
		return // Skip nil statements
	}
	switch s := stmt.(type) {
	case *parser.PrintStatement:
		if s == nil || s.Value == nil {
			if s != nil {
				fmt.Printf("Error at line %d, col %d: Invalid print statement \n",
					s.Token().Line, s.Token().Column)
			}
			return
		}
		value := i.evalExpression(s.Value)
		if value != nil {
			fmt.Println(value.String())
		} else {
			fmt.Printf("Error at line %d, col %d: Invalid expression in print \n",
				s.Token().Line, s.Token().Column)
		}
	case *parser.AssignmentStatement:
		if s == nil || s.Name == nil || s.Value == nil {
			if s != nil {
				fmt.Printf("Error at line %d, col %d: Invalid assignment statement \n",
					s.Token().Line, s.Token().Column)
			}
			return
		}
		value := i.evalExpression(s.Value)
		if value != nil {
			i.env.store[s.Name.Value] = value
		} else {
			fmt.Printf("Error at line %d, col %d: Invalid expression in assignment \n",
				s.Token().Line, s.Token().Column)
		}
	case *parser.IfStatement:
		if s == nil || s.Condition == nil {
			if s != nil {
				fmt.Printf("Error at line %d, col %d: Invalid if statement \n",
					s.Token().Line, s.Token().Column)
			}
			return
		}
		if s.Consequence == nil {
			fmt.Printf("Error at line %d, col %d: Invalid if block \n",
				s.Token().Line, s.Token().Column)
			return
		}
		condition := i.evalExpression(s.Condition)
		if condition == nil {
			fmt.Printf("Error at line %d, col %d: Invalid condition in if \n",
				s.Token().Line, s.Token().Column)
			return
		}
		if isTruthy(condition) {
			for _, stmt := range s.Consequence.Statements {
				if stmt != nil {
					i.evalStatement(stmt)
				}
			}
		} else if s.Alternative != nil && !isTruthy(condition) {
			for _, stmt := range s.Alternative.Statements {
				if stmt != nil {
					i.evalStatement(stmt)
				}
			}
		}
	default:
		// Handle cases where we can't get token info
		fmt.Printf("Error: Unknown statement type\n")
	}
}

// evalExpression evaluates an expression and returns an Object.
func (i *Interpreter) evalExpression(expr parser.Expression) Object {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *parser.NumberLiteral:
		return &IntObject{Value: e.Value}
	case *parser.StringLiteral:
		return &StringObject{Value: e.Value}
	case *parser.Identifier:
		value, ok := i.env.store[e.Value]
		if !ok {
			fmt.Printf("Error at line %d, col %d: Undefined variable %s \n",
				e.Token.Line, e.Token.Column, e.Value)
			return nil
		}
		return value
	case *parser.BinaryExpression:
		left := i.evalExpression(e.Left)
		if left == nil {
			return nil
		}
		right := i.evalExpression(e.Right)
		if right == nil {
			return nil
		}
		return i.evalBinaryExpression(e.Token, left, e.Operator, right)
	default:
		// Try to get token info if possible, else use -1
		line, col := -1, -1
		if tokExpr, ok := expr.(interface{ Token() lexer.Token }); ok {
			tok := tokExpr.Token()
			line, col = tok.Line, tok.Column
		}
		fmt.Printf("Error at line %d, col %d: Unknown expression type \n",
			line, col)
		return nil
	}
}

// evalBinaryExpression evaluates a binary expression (arithmetic or comparison).
func (i *Interpreter) evalBinaryExpression(token lexer.Token, left Object, op string, right Object) Object {
	// Handle arithmetic (int + int)
	if leftInt, ok1 := left.(*IntObject); ok1 {
		if rightInt, ok2 := right.(*IntObject); ok2 {
			switch op {
			case "+":
				return &IntObject{Value: leftInt.Value + rightInt.Value}
			case "-":
				return &IntObject{Value: leftInt.Value - rightInt.Value}
			case "*":
				return &IntObject{Value: leftInt.Value * rightInt.Value}
			case "%":
				return &IntObject{Value: leftInt.Value % rightInt.Value}
			case "/":
				if rightInt.Value == 0 {
					fmt.Printf("Error at line %d, col %d: Division by zero \n",
						token.Line, token.Column)
					return nil
				}
				return &IntObject{Value: leftInt.Value / rightInt.Value}
			case "==":
				return &IntObject{Value: boolToInt(leftInt.Value == rightInt.Value)}
			case "!=":
				return &IntObject{Value: boolToInt(leftInt.Value != rightInt.Value)}
			case "<":
				return &IntObject{Value: boolToInt(leftInt.Value < rightInt.Value)}
			case ">":
				return &IntObject{Value: boolToInt(leftInt.Value > rightInt.Value)}
			case "<=":
				return &IntObject{Value: boolToInt(leftInt.Value <= rightInt.Value)}
			case ">=":
				return &IntObject{Value: boolToInt(leftInt.Value >= rightInt.Value)}
			}
		}
	}
	// Handle string + string (concatenation)
	if leftStr, ok1 := left.(*StringObject); ok1 {
		if rightStr, ok2 := right.(*StringObject); ok2 {
			if op == "+" {
				return &StringObject{Value: leftStr.Value + rightStr.Value}
			}
		}
	}
	fmt.Printf("Error at line %d, col %d: Invalid operation %s between %s and %s \n",
		token.Line, token.Column, op, left.String(), right.String())
	return nil
}

// boolToInt converts a boolean to 1 (true) or 0 (false).
func boolToInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

// isTruthy determines if an Object is truthy for conditionals.
func isTruthy(obj Object) bool {
	switch o := obj.(type) {
	case *IntObject:
		return o.Value != 0
	case *StringObject:
		return len(o.Value) > 0
	default:
		return false
	}
}
