package parser

import (
	"fmt"
	"strconv"

	"github.com/salillakra/npp/frontend/lexer"
)

// Node represents a node in the AST.
type Node interface {
	String() string
}

// Statement represents a statement node.
type Statement interface {
	Node
	statementNode()
	Token() lexer.Token
}

// Expression represents an expression node.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out string
	for _, stmt := range p.Statements {
		if stmt != nil {
			out += stmt.String() + "\n"
		}
	}
	return out
}

// PrintStatement represents a print statement (e.g., suna "You suck!").
type PrintStatement struct {
	Tok   lexer.Token
	Value Expression
}

func (ps *PrintStatement) statementNode()     {}
func (ps *PrintStatement) String() string     { return fmt.Sprintf("suna %s", ps.Value.String()) }
func (ps *PrintStatement) Token() lexer.Token { return ps.Tok }

// AssignmentStatement represents an assignment statement (e.g., sun x = 69).
type AssignmentStatement struct {
	Tok   lexer.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("sun %s = %s", as.Name.String(), as.Value.String())
}
func (as *AssignmentStatement) Token() lexer.Token { return as.Tok }

// IfStatement represents an if statement (e.g., agar x > 50 { ... }).
type IfStatement struct {
	Tok         lexer.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	if is.Alternative != nil {
		return fmt.Sprintf("agar %s { ... } magar { ... }", is.Condition.String())
	}
	return fmt.Sprintf("agar %s { ... }", is.Condition.String())
}
func (is *IfStatement) Token() lexer.Token { return is.Tok }

// BlockStatement represents a block of statements (e.g., { suna 42; }).
type BlockStatement struct {
	Tok        lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()     {}
func (bs *BlockStatement) String() string     { return "{ ... }" }
func (bs *BlockStatement) Token() lexer.Token { return bs.Tok }

// Identifier represents a variable name (e.g., x).
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

// NumberLiteral represents an integer literal (e.g., 69).
type NumberLiteral struct {
	Token lexer.Token
	Value int64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return fmt.Sprintf("%d", nl.Value) }

// StringLiteral represents a string literal (e.g., "You suck!").
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("%q", sl.Value) }

// BinaryExpression represents a binary operation (e.g., x + 10, x > y).
type BinaryExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}

// Parser holds the lexer and current/peek tokens.
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	Debug     bool
}

// New creates a new Parser.
func New(l *lexer.Lexer, Debug bool) *Parser {
	p := &Parser{l: l, Debug: Debug}
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances to the next token.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses the input into a Program AST.
func (p *Parser) ParseProgram() *Program {
	program := &Program{Statements: []Statement{}}
	for p.curToken.Type != lexer.EOF {
		if p.Debug {
			fmt.Printf("Debug: Parsing statement at %v (line %d, col %d)\n", p.curToken, p.curToken.Line, p.curToken.Column)

		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			fmt.Printf("Error at line %d, col %d: Invalid statement, got %s // Keep it together, genius!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
			p.nextToken()
		}
		// Skip optional semicolons
		for p.curToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
	}
	return program
}

// parseStatement parses a single statement.
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case lexer.SUN:
		// Expect: sun IDENT = expression
		stmt := &AssignmentStatement{Tok: p.curToken}
		p.nextToken()
		if p.curToken.Type != lexer.IDENT {
			fmt.Printf("Error at line %d, col %d: Expected identifier after SUN, got %s // My grandma codes better!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
			return nil
		}
		stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		if p.curToken.Type != lexer.ASSIGN {
			fmt.Printf("Error at line %d, col %d: Expected = after identifier, got %s // Yo, nice one, jerk!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
			return nil
		}
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
		if stmt.Value == nil {
			fmt.Printf("Error at line %d, col %d: Expected expression after =, got %s // You absolute walnut!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
			return nil
		}
		return stmt
	case lexer.SUNA:
		return p.parsePrintStatement()
	case lexer.AGAR:
		return p.parseIfStatement()
	default:
		fmt.Printf("Error at line %d, col %d: Invalid statement, got %s // Keep it together, genius!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
		return nil
	}
}

// parsePrintStatement parses a print statement (e.g., suna "You suck!" or suna x).
func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{Tok: p.curToken}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil {
		fmt.Printf("Error at line %d, col %d: Expected expression after suna, got %s // You absolute walnut!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
		return nil
	}
	return stmt
}

// parseIfStatement parses an if statement (e.g., agar x > 50 { ... } magar { ... }).
func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Tok: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if stmt.Condition == nil {
		fmt.Printf("Error at line %d, col %d: Expected condition after agar, got %s // This syntax sucks, fix it!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
		return nil
	}
	if p.curToken.Type != lexer.LBRACE {
		fmt.Printf("Error at line %d, col %d: Expected { after condition, got %s // Get your braces together, loser!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
		return nil
	}
	stmt.Consequence = p.parseBlockStatement()
	if stmt.Consequence == nil {
		fmt.Printf("Error at line %d, col %d: Invalid block after agar // This ain't working, jerk!\n", p.curToken.Line, p.curToken.Column)
		return nil
	}
	p.nextToken() // Skip closing brace
	for p.curToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}
	if p.curToken.Type == lexer.MAGAR {
		p.nextToken()
		if p.curToken.Type != lexer.LBRACE {
			fmt.Printf("Error at line %d, col %d: Expected { after magar, got %s // Get your braces together, loser!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
			return nil
		}
		stmt.Alternative = p.parseBlockStatement()
		if stmt.Alternative == nil {
			fmt.Printf("Error at line %d, col %d: Invalid block after magar // This ain't working, jerk!\n", p.curToken.Line, p.curToken.Column)
			return nil
		}
		p.nextToken() // Skip closing brace
	}
	return stmt
}

// parseBlockStatement parses a block of statements (e.g., { suna 42; }).
func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Tok: p.curToken, Statements: []Statement{}}
	p.nextToken()
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		} else {
			p.nextToken()
		}
		for p.curToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
	}
	if p.curToken.Type != lexer.RBRACE {
		fmt.Printf("Error at line %d, col %d: Expected } to close block, got %s // Close your blocks, you walnut!\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
		return nil
	}
	return block
}

// Precedence levels for operators
const (
	LOWEST      = 1
	EQUALS      = 2 // ==, !=
	LESSGREATER = 3 // <, >, <=, >=
	SUM         = 4 // +, -
	PRODUCT     = 5 // *, /
)

var precedences = map[lexer.TokenType]int{
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.LE:       LESSGREATER,
	lexer.GE:       LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTERISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
}

// parseExpression parses an expression with precedence handling.
func (p *Parser) parseExpression(precedence int) Expression {
	var left Expression
	if p.curToken.Type == lexer.MINUS {
		token := p.curToken
		p.nextToken()
		if p.curToken.Type != lexer.INT {
			fmt.Printf("Error at line %d, col %d: Expected number after -, got %s // Numbers too hard for you, huh?\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
			return nil
		}
		value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
		if err != nil {
			fmt.Printf("Error at line %d, col %d: Invalid number %s // Numbers too hard for you, huh?\n", p.curToken.Line, p.curToken.Column, p.curToken.Literal)
			return nil
		}
		left = &NumberLiteral{Token: token, Value: -value}
		p.nextToken()
	} else {
		left = p.parsePrimary()
		if left == nil {
			return nil
		}
	}

	for p.curToken.Type != lexer.EOF &&
		p.curToken.Type != lexer.SEMICOLON &&
		p.curToken.Type != lexer.RBRACE &&
		p.curToken.Type != lexer.LBRACE &&
		precedence < p.getCurrentPrecedence() {
		if !isOperator(p.curToken.Type) {
			break
		}
		op := p.curToken
		p.nextToken()
		right := p.parseExpression(p.getPrecedence(op.Type))
		if right == nil {
			fmt.Printf("Error at line %d, col %d: Expected expression after %s // What's this nonsense, loser?\n", p.curToken.Line, p.curToken.Column, op.Literal)
			return nil
		}
		left = &BinaryExpression{Token: op, Left: left, Operator: op.Literal, Right: right}
	}
	return left
}

// parsePrimary parses a primary expression (number, string, or identifier).
func (p *Parser) parsePrimary() Expression {
	switch p.curToken.Type {
	case lexer.INT:
		value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
		if err != nil {
			fmt.Printf("Error at line %d, col %d: Invalid number %s // Numbers too hard for you, huh?\n", p.curToken.Line, p.curToken.Column, p.curToken.Literal)
			return nil
		}
		result := &NumberLiteral{Token: p.curToken, Value: value}
		p.nextToken()
		return result
	case lexer.STRING:
		result := &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		return result
	case lexer.IDENT:
		result := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		return result
	default:
		fmt.Printf("Error at line %d, col %d: Expected number, string, or identifier, got %s // What even is this, genius?\n", p.curToken.Line, p.curToken.Column, p.curToken.Type)
		return nil
	}
}

// getCurrentPrecedence returns the precedence of the current token.
func (p *Parser) getCurrentPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// getPrecedence returns the precedence of the given token type.
func (p *Parser) getPrecedence(tokenType lexer.TokenType) int {
	if p, ok := precedences[tokenType]; ok {
		return p
	}
	return LOWEST
}

// isOperator checks if a token is an operator.
func isOperator(tokenType lexer.TokenType) bool {
	return tokenType == lexer.PLUS || tokenType == lexer.MINUS ||
		tokenType == lexer.ASTERISK || tokenType == lexer.SLASH ||
		tokenType == lexer.EQ || tokenType == lexer.NOT_EQ ||
		tokenType == lexer.LT || tokenType == lexer.GT ||
		tokenType == lexer.LE || tokenType == lexer.GE
}
