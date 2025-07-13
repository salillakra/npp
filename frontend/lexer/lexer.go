package lexer

import (
	"unicode"
)

type Lexer struct {
	input        string // source code
	position     int    // current position (index of ch)
	readPosition int    // position after current char
	ch           byte   // current char
	line         int    // current line number (1-based)
	column       int    // current column number (1-based)
}

// New creates a new Lexer instance for the given input string.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 1}
	l.readChar()
	return l
}

// readChar advances the lexer to the next character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else if l.ch != 0 {
		l.column++
	}
}

// TokenType represents the type of a token.
type TokenType string

// Token represents a lexical token with its type, literal value, and position.
type Token struct {
	Type    TokenType // Token type (e.g., SUN, INT)
	Literal string    // Literal value (e.g., "69", "x")
	Line    int       // Line number (1-based)
	Column  int       // Column number (1-based)
}

// Token types
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT  = "IDENT"  // x, y, jerk
	INT    = "INT"    // 123
	STRING = "STRING" // "you suck"

	// Operators
	ASSIGN   = "="
	EQ       = "=="
	NOT_EQ   = "!="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	LE       = "<="
	GE       = ">="

	// Punctuation
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	SUN   = "SUN"   // sun (variable declaration)
	SUNA  = "SUNA"  // suna (print)
	AGAR  = "AGAR"  // agar (if)
	MAGAR = "MAGAR" // magar (else)
	GLOW  = "GLOW"  // glow (function)
	FHEK  = "FHEK"  // fhek (return)
	YAS   = "YAS"   // yas (true)
	NAH   = "NAH"   // nah (false)
	GRIND = "GRIND" // grind (while)
)

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	l.skipComment()

	tok := Token{Line: l.line, Column: l.column}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(ASSIGN, string(l.ch), l.line, l.column)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(BANG, string(l.ch), l.line, l.column)
		}
	case '+':
		tok = newToken(PLUS, string(l.ch), l.line, l.column)
	case '-':
		tok = newToken(MINUS, string(l.ch), l.line, l.column)
	case '*':
		tok = newToken(ASTERISK, string(l.ch), l.line, l.column)
	case '/':
		if l.peekChar() == '/' {
			l.readChar() // Skip first '/'
			l.skipComment()
			return l.NextToken() // Recursively get next token after comment
		}
		tok = newToken(SLASH, string(l.ch), l.line, l.column)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(LT, string(l.ch), l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(GT, string(l.ch), l.line, l.column)
		}
	case ',':
		tok = newToken(COMMA, string(l.ch), l.line, l.column)
	case ';':
		tok = newToken(SEMICOLON, string(l.ch), l.line, l.column)
	case '(':
		tok = newToken(LPAREN, string(l.ch), l.line, l.column)
	case ')':
		tok = newToken(RPAREN, string(l.ch), l.line, l.column)
	case '{':
		tok = newToken(LBRACE, string(l.ch), l.line, l.column)
	case '}':
		tok = newToken(RBRACE, string(l.ch), l.line, l.column)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = l.line
		tok.Column = l.column
		return tok
	case 0:
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = newToken(ILLEGAL, string(l.ch), l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token with the given type, literal, and position.
func newToken(tokenType TokenType, literal string, line, column int) Token {
	return Token{Type: tokenType, Literal: literal, Line: line, Column: column}
}

// skipWhitespace skips over whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment skips single-line comments starting with "//".
func (l *Lexer) skipComment() {
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

// readIdentifier reads an identifier or keyword.
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readNumber reads an integer literal.
func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readString reads a string literal enclosed in quotes.
func (l *Lexer) readString() string {
	l.readChar() // Skip opening quote
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	if l.ch == 0 {
		return l.input[start:l.position] // Unterminated string
	}
	str := l.input[start:l.position]
	l.readChar() // Skip closing quote
	return str
}

// peekChar returns the next character without advancing the lexer.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// isLetter checks if a character is a letter or underscore.
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

// lookupIdent maps identifiers to keyword token types.
func lookupIdent(ident string) TokenType {
	keywords := map[string]TokenType{
		"sun":   SUN,
		"suna":  SUNA,
		"agar":  AGAR,
		"magar": MAGAR,
		"glow":  GLOW,
		"fhek":  FHEK,
		"yas":   YAS,
		"nah":   NAH,
		"grind": GRIND,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
