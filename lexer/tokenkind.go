package lexer

type TokenKind int

const (
	EOF TokenKind = iota
	EOL

	// Keywords
	MOD
	USE
	ALIAS
	IF
	ELSE
	STRUCT
	ENUM
	VAR
	CONST
	RDO
	REF
	ATTR

	// Symbols
	PLUS
	PLUSEQ
	MINUS
	MINUSEQ
	STAR
	STAREQ
	SLASH
	SLASHEQ
	PERCENT
	PERCENTEQ

	LPAREN  // (
	RPAREN  // )
	LBRACE  // {
	RBRACE  // }
	LSQUARE // [
	RSQUARE // ]
	DOLLAR  // $
	COMMA   // ,
	BANG    // !

	DOT        // .
	DOTDOT     // ..
	COLON      // :
	COLONCOLON // ::
	EQ         // =
	EQEQ       // ==

	// DYN
	NAME // [_a-zA-Z][_a-zA-Z0-9]*
	STRING
	NUMBER  // Whole number
	DECIMAL // Decimal number
)
