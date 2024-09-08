package lexer

import (
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/golang-collections/collections/stack"
)

const tab_depth = 4 // tabs are equivilent to four spaces, but tabs should not be used

type Lexer struct {
	stream      *CharStream
	scopeDepths *stack.Stack // Stack of ints
	scopes      *stack.Stack // Stack of *Lines
}

func NewLexer(s *CharStream) Lexer {
	return Lexer{s, stack.New(), stack.New()}
}

func (l *Lexer) Parse() []*Line {
	lines := make([]*Line, 0, 20) // just preemptive capacity, all top level things
	l.stream.Next()
	var prevLine *Line
	var prevDepth int

	for {
		depth := l.countAndSkipIndent()
		line, eof := l.parseLine()
		// fmt.Println(123)

		if eof {
			break
		}

		if len(line.Content) == 0 {
			continue
		}
		if prevLine == nil && depth != 0 {
			log.Fatalln("Program cannot start with indentation")
		}

		if depth > prevDepth {
			l.scopes.Push(prevLine)
			l.scopeDepths.Push(prevDepth)
		}
		if depth < prevDepth {
			for {

				topDepth := l.scopeDepths.Peek()

				if d, ok := topDepth.(int); ok {
					if d == depth || d == 0 {
						break
					}
				} else {
					log.Fatalf("Unknown indentation amount '%d' at line %d\n", depth, l.stream.lineIdx)
				}

				// fmt.Println(l.scopeDepths.Pop())
				l.scopes.Pop()
			}
		}

		if depth == 0 {
			lines = append(lines, line)
		} else {
			l.scopes.Peek().(*Line).Children = append(l.scopes.Peek().(*Line).Children, line)
		}

		prevDepth = depth
		prevLine = line
	}
	// fmt.Println(lines)
	return lines
}

func (l *Lexer) countAndSkipIndent() (depth int) {
	prevDepth := 0
	for {
		c := l.stream.Current()

		if c == ' ' {
			depth++
		}
		if c == '\t' {
			depth += tab_depth
		}

		if prevDepth == depth {
			break
		}

		l.stream.Next()

		prevDepth = depth
	}
	return
}

// Assumes charstream is at the beginning of a line, and that the preceding whitespace is skipped
func (l *Lexer) parseLine() (line *Line, isEOF bool) {
	toks := make([]Token, 0, 5)
	isEOF = false

	for {
		tok := l.parseToken()
		if tok.Kind == EOL {
			break
		}
		if tok.Kind == EOF {
			isEOF = true
			break
		}

		toks = append(toks, tok)
	}

	line = &Line{toks, make([]*Line, 0)}
	return
}

// Not for used to manage or count indentation
func (l *Lexer) skipWs() {
	c := l.stream.Current()
	for c == ' ' || c == '\t' || c == '\r' {
		c = l.stream.Next()
	}
}

func (l *Lexer) parseToken() Token {
	//#region EOF and EOL
	if l.stream.Current() == '\n' {
		l.stream.Next()
		return Token{EOL, nil}
	}
	if l.stream.Current() == 0 {
		return Token{EOF, nil}
	}
	l.skipWs()
	c := l.stream.Current()
	if c == 0 {
		return Token{EOF, nil}
	}
	if c == '\n' {
		l.stream.Next()
		return Token{EOL, nil}
	}

	//#endregion EOF and EOL

	c = l.stream.Current()

	//#region SYMBOLS
	if tok, ok := l.matchSymbols(map[rune]TokenKind{
		'(': LPAREN,
		')': RPAREN,
		'{': RBRACE,
		'}': RBRACE,
		'[': LSQUARE,
		']': RSQUARE,
		'$': DOLLAR,
		',': COMMA,
		'!': BANG,
		'&': AMPERSAND,
	}); ok {
		return Token{tok, nil}
	}

	if tok, ok := l.singleOrDoubleToken('.', DOT, DOTDOT); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.singleOrDoubleToken(':', COLON, COLONCOLON); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.singleOrDoubleToken('=', EQ, EQEQ); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.singleOrDoubleToken('<', LESSTHAN, SHIFTLEFT); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.singleOrDoubleToken('<', GREATERTHAN, SHIFTRIGHT); ok {
		return Token{tok, nil}
	}

	if tok, ok := l.multiToken('+', PLUS, map[rune]TokenKind{'=': PLUSEQ}); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.multiToken('-', MINUS, map[rune]TokenKind{'=': MINUSEQ}); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.multiToken('*', STAR, map[rune]TokenKind{'=': STAREQ}); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.multiToken('/', SLASH, map[rune]TokenKind{'=': SLASHEQ}); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.multiToken('%', PERCENT, map[rune]TokenKind{'=': PERCENTEQ}); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.multiToken('<', LESSTHAN, map[rune]TokenKind{'<': SHIFTLEFT, '=': LESSEQ}); ok {
		return Token{tok, nil}
	}
	if tok, ok := l.multiToken('>', GREATERTHAN, map[rune]TokenKind{'>': SHIFTRIGHT, '=': GREATEREQ}); ok {
		return Token{tok, nil}
	}

	//#endregion SYMBOLS

	//#region Advanced symbols
	// ToDo, add mul and div
	wasArithEquable := false
	var arithTok TokenKind = EOF
	if c == '+' {
		wasArithEquable = true
		arithTok = PLUS
	}
	if c == '-' {
		wasArithEquable = true
		arithTok = MINUS
	}
	if c == '*' {
		wasArithEquable = true
		arithTok = STAR
	}
	if c == '/' {
		wasArithEquable = true
		arithTok = SLASH
	}
	if c == '%' {
		wasArithEquable = true
		arithTok = PERCENT
	}
	if wasArithEquable {
		if l.stream.Next() == '=' {
			arithTok = TokenKind(int(arithTok) + 1) // Hacky shit but idc it works
			l.stream.Next()
		}
		if l.stream.Current() == '/' {
			// Skip comment
			for {
				c := l.stream.Next()
				if c == '\n' || c == eof_char {
					break
				}
			}
			return l.parseToken()
		}

		return Token{arithTok, nil}
	}
	//#endregion Advanced symbols

	//#region Litterals
	if unicode.IsLetter(c) || c == '_' {
		sb := strings.Builder{}
		sb.Reset()
		sb.WriteRune(c)

		// fmt.Printf("First %c\n", c)
		for unicode.IsLetter(l.stream.Next()) ||
			unicode.IsDigit(l.stream.Current()) ||
			l.stream.Current() == '_' {

			sb.WriteRune(l.stream.Current())
		}

		str := sb.String()

		switch str {
		case "mod":
			return Token{MOD, nil}
		case "use":
			return Token{USE, nil}
		case "alias":
			return Token{ALIAS, nil}
		case "if":
			return Token{IF, nil}
		case "else":
			return Token{ELSE, nil}
		case "struct":
			return Token{STRUCT, nil}
		case "enum":
			return Token{ENUM, nil}
		case "var":
			return Token{VAR, nil}
		case "const":
			return Token{CONST, nil}
		case "rdo":
			return Token{RDO, nil}
		case "ref":
			return Token{REF, nil}
		case "attr":
			return Token{ATTR, nil}
		case "async":
			return Token{ASYNC, nil}
		case "do":
			return Token{DO, nil}
		case "for":
			return Token{FOR, nil}
		case "return":
			return Token{RETURN, nil}
		}
		return Token{NAME, str}
	}

	// handle strings someday
	if c == '"' {
		sb := strings.Builder{}

		for {
			next := l.stream.Next()
			if next == '"' {
				break
			}
			if next == '\n' || next == eof_char {
				log.Fatalf("String was unterminated at line %d\n", l.stream.lineIdx)
			}
			sb.WriteRune(next)
		}

		l.stream.Next()
		return Token{Kind: STRING, Data: sb.String()}
	}

	if unicode.IsDigit(c) {
		sb := strings.Builder{}
		sb.WriteRune(c)
		isDecimal := false
		for {
			c := l.stream.Next()

			if unicode.IsDigit(c) {
				sb.WriteRune(c)
			} else if c == '.' {
				if isDecimal {
					break // using dot notation on a float
				}
				isDecimal = true
				sb.WriteRune(c)
			} else if unicode.IsLetter(c) {
				// ToDo: handle special numbers
				log.Fatalf("Numbers cannot be followed by letters at line %d\n", l.stream.lineIdx+1)
			} else {
				break
			}
		}
		if isDecimal {
			f, err := strconv.ParseFloat(sb.String(), 64)
			if err != nil {
				panic(err) // unreachable
			}
			return Token{DECIMAL, f}
		} else {
			i, err := strconv.ParseInt(sb.String(), 10, 64)
			if err != nil {
				panic(err) // unreachable
			}
			return Token{NUMBER, i}
		}
	}
	//#endregion Litterals

	log.Fatalf("Unexpected char %c at line %d", c, l.stream.lineIdx)
	return Token{}
}

func (l *Lexer) singleOrDoubleToken(c rune, single, double TokenKind) (TokenKind, bool) {
	if l.stream.Current() != c {
		return 0, false
	}

	if l.stream.Next() == c {
		l.stream.Next()
		return double, true
	}

	return single, true
}

func (l *Lexer) multiToken(c rune, single TokenKind, extra map[rune]TokenKind) (TokenKind, bool) {
	if l.stream.Current() != c {
		return 0, false
	}

	current := l.stream.Next()

	for ex, t := range extra {
		if current == ex {
			l.stream.Next()
			return t, true
		}
	}

	return single, true
}

func (l *Lexer) matchSymbols(toks map[rune]TokenKind) (TokenKind, bool) {
	c := l.stream.Current()
	for key, value := range toks {
		if c == key {
			l.stream.Next()
			return value, true
		}
	}
	return 0, false
}
