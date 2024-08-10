package lexer

type TokenKind int

const (
	EOF TokenKind = iota
	EOL

	// Keywords
	MOD
	IF
	ELSE

	// Symbols
	PLUS
	PLUSEQ
	MINUS
	MINUSEQ

	DOT
	LPAREN
	RPAREN

	// DYN
	NAME // [_a-zA-Z][_a-zA-Z0-9]*
	STRING
	NUMBER  // Whole number
	DECIMAL // Decimal number
)
