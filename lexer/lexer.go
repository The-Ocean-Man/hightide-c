package lexer

import (
	"log"
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

func (l *Lexer) depthGreater(depth int) bool {
	top := l.scopeDepths.Peek()
	if top == nil {
		return depth > 0
	}
	return depth > top.(int)
}

func (l *Lexer) depthLess(depth int) bool {
	top := l.scopeDepths.Peek()
	if top == nil {
		return false
	}
	return depth < top.(int)
}

func (l *Lexer) Parse() []*Line {
	lines := make([]*Line, 0, 20) // just preemptive capacity, all top level things
	l.stream.Next()
	var prevLine *Line
	var prevDepth int

	for {
		depth := l.countAndSkipIndent()
		line, eof := l.parseLine()

		if l.depthGreater(depth) {
			l.scopes.Push(prevLine)
			l.scopeDepths.Push(prevDepth)
		}
		if l.depthLess(depth) {
			for {
				prevDepthPtr := l.scopeDepths.Peek()
				var prevDepth int = 0
				if prevDepthPtr != nil {
					prevDepth = prevDepthPtr.(int)
				}
				if prevDepth == depth {
					break
				}

				l.scopes.Pop()
				if l.scopeDepths.Pop() == nil {
					log.Fatalln("Scopes broke in lexer")
				}
			}
		}

		if top := l.scopes.Peek(); top != nil {
			top.(*Line).Children = append(top.(*Line).Children, line)
		} else {
			lines = append(lines, line)
		}

		if eof {
			// fmt.Println("dead")
			break
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

	c = l.stream.Current()

	if c == '(' {
		l.stream.Next()
		return Token{LPAREN, nil}
	} else if c == ')' {
		l.stream.Next()
		return Token{RPAREN, nil}
	} else if c == '.' {
		l.stream.Next()
		return Token{DOT, nil}
	}

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
	if wasArithEquable {
		if l.stream.Next() == '=' {
			arithTok = TokenKind(int(arithTok) + 1) // Hacky shit but idc
			l.stream.Next()
		}
		return Token{arithTok, nil}
	}

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
		case "if":
			return Token{IF, nil}
		case "else":
			return Token{ELSE, nil}
		}
		return Token{NAME, str}
	}

	// handle strings someday

	log.Fatalf("Unexpected char %c at line %d", c, l.stream.lineIdx)
	return Token{}
}
