package parser

import (
	"github.com/The-Ocean-Man/hightide-c/lexer"
)

type LineWalker struct {
	ln   *lexer.Line
	curr lexer.Token
	idx  int
}

func MakeWalker(l *lexer.Line) *LineWalker {
	w := LineWalker{l, lexer.Token{}, 0}
	w.Next()
	return &w
}

func (w *LineWalker) Next() lexer.Token {
	if w.idx >= len(w.ln.Content) {
		w.curr = lexer.Token{Kind: lexer.EOL, Data: nil}
	} else {
		w.curr = w.ln.Content[w.idx]
	}
	w.idx++

	return w.curr
}

func (w *LineWalker) Get() lexer.Token {
	return w.curr
}

func (w *LineWalker) GetLine() *lexer.Line {
	return w.ln
}
