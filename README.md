# NPP: Nasty Plus Plus

NPP is a custom programming language implemented in Go. It features a simple interpreter, lexer, and parser, supporting variable declarations, arithmetic, string operations, conditionals, and print statements. This project demonstrates how to build a basic language and interpreter from scratch.

## Features

- Variable declaration and assignment (numbers and strings)
- Arithmetic and string concatenation
- Print statements
- Conditional statements (`agar`/`magar`)

## Project Structure

```
core/
  interpreter/         # Interpreter logic
frontend/
  lexer/               # Lexical analyzer
  parser/              # Parser and AST
main/
  main.go              # Entry point for running NPP code
  hello.npp            # Example NPP program
  test_test.go         # Unit tests
```

## Example NPP Program (`main/hello.npp`)

```
sun number = 2;
suna number;
sun a = "hello ";
sun b = "world";
suna a+b;
agar number >= 10 {
    suna "nice one baby!";
} magar{
    suna "fuck off!";
}
```

## Usage

### 1. Build and Run

```sh
cd main
# Run the interpreter with an NPP file
# Usage: go run main.go <file.npp>
go run main.go hello.npp
```

### 2. Run Tests

```sh
cd main
go test
```

## Language Reference

- `sun <var> = <value>;` — Declare and assign a variable
- `suna <expr>;` — Print an expression
- `agar <condition> { ... } magar { ... }` — If/else conditional
- Supports `+`, `-`, `*`, `/`, `%`, `==`, `!=`, `<`, `>`, `<=`, `>=` operators

## Development

- Written in Go
- Modular structure for easy extension
- Add new statements or expressions by editing the parser and interpreter

## Contributing

Pull requests and issues are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

MIT License
