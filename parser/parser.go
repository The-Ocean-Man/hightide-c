package parser

import (
	"fmt"
	"log"

	"github.com/The-Ocean-Man/hightide-c/ast"
	lex "github.com/The-Ocean-Man/hightide-c/lexer"
)

type Parser struct {
	lines  []*lex.Line
	walker *LineWalker
}

func MakeParser(l []*lex.Line) Parser {
	return Parser{l, nil}
}

func (p *Parser) ParseProgram() *ast.ProgramNode {
	topLevelStatements := make([]ast.Node, 0)

	for _, ln := range p.lines {
		walker := MakeWalker(ln)
		p.walker = &walker
		if isVarDecToken(walker.Get().Kind) {
			topLevelStatements = append(topLevelStatements, p.ParseVarDec())
		} else if isExprStartToken(walker.Get().Kind) {
			topLevelStatements = append(topLevelStatements, p.ParseExpr())
		}
	}

	return &ast.ProgramNode{Children: topLevelStatements}
}

// Returns true if kind == VAR, CONST or RDO
func isVarDecToken(kind lex.TokenKind) bool {
	return kind == lex.VAR || kind == lex.CONST || kind == lex.RDO
}

func isExprStartToken(kind lex.TokenKind) bool {
	return kind == lex.NUMBER || kind == lex.DECIMAL || kind == lex.LPAREN || kind == lex.MINUS || kind == lex.BANG // kind == lex.
}

func (p *Parser) ParseVarDec() *ast.VarDecNode {
	decTokKind := p.ExpectAnyOf(lex.VAR, lex.CONST, lex.RDO).Kind
	_ = decTokKind
	var mut ast.Mutability = ast.Mutable
	switch decTokKind {
	case lex.CONST:
		mut = ast.Constant
	case lex.RDO:
		mut = ast.ReadOnly
	case lex.VAR:
		mut = ast.Mutable
	}
	varName := p.Expect(lex.NAME).Data.(string)
	return &ast.VarDecNode{Name: varName, Mut: mut, Type: nil, Value: nil}
}

// Oh boy this is gonna be long
func (p *Parser) ParseExpr() ast.Node {
	bits := make([]ExprBit, 0, 3) // three cap cuz thats usually the minimum size

	for {
		bit := p.ParseExprBit()
		if bit == nil {
			break
		}
		bits = append(bits, bit)
	}

	return ParseExpression(bits)
}
func (p *Parser) ParseExprBit() ExprBit {
	term := p.ParseExprTerminal()

	// account for invocations and indexing
	return term
}

func (p *Parser) ParseExprTerminal() ExprBit {
	// current := p.walker.Get()
	if p.IsCurrent(lex.EOL, lex.EOF, lex.RPAREN, lex.COMMA) { // end of expr
		return nil
	}
	if p.IsCurrent(lex.STRING) {
		strTok := p.Expect(lex.STRING)
		return &ast.StringLitteralNode{Value: strTok.Data.(string)}
	}
	if p.IsCurrent(lex.NUMBER) {
		intTok := p.Expect(lex.NUMBER)
		return &ast.IntLitteralNode{Value: intTok.Data.(int64)}
	}
	if p.IsCurrent(lex.DECIMAL) {
		floatTok := p.Expect(lex.DECIMAL)
		return &ast.FloatLitteralNode{Value: floatTok.Data.(float64)}
	}

	if p.IsCurrent(lex.NAME) {
		return p.ParseIdent()
	}
	if p.IsCurrent(lex.LPAREN) {
		fmt.Println("paren")
		return p.ParseParenExpr()
	}
	if t, ok := p.Optional(lex.PLUS, lex.MINUS, lex.STAR, lex.SLASH, lex.PERCENT, lex.BANG,
		lex.REF, lex.CONST, lex.RDO, lex.DOLLAR); ok {
		return t.Kind
	}

	// else
	log.Fatalf("Unexpected token %d in expr", p.walker.Get().Kind)

	return nil // unreachable
}

func (p *Parser) ParseParenExpr() ast.Node {
	p.Expect(lex.LPAREN)
	e := p.ParseExpr()
	p.Expect(lex.RPAREN)
	return e
}

func (p *Parser) ParseIdent() *ast.IdentNode {
	first := p.Expect(lex.NAME)
	top := &ast.IdentNode{}
	top.Name = first.Data.(string)
	node := top

	for {
		if _, ok := p.Optional(lex.DOT); ok {
			node.UsedDot = true
		} else if _, ok := p.Optional(lex.COLONCOLON); ok {
			node.UsedDot = false
		} else {
			break
		}

		name := p.Expect(lex.NAME)
		lower := &ast.IdentNode{}
		lower.Name = name.Data.(string)
		node.Child = lower
		node = lower
	}

	return top
}

func (p *Parser) Expect(kind lex.TokenKind) lex.Token {
	t := p.walker.Get()
	if t.Kind == kind {
		p.walker.Next()
		return t
	}
	log.Fatalf("Expected %d but got %d instead", kind, t.Kind)
	return lex.Token{}
}

func (p *Parser) ExpectAnyOf(kinds ...lex.TokenKind) lex.Token {
	t := p.walker.Get()
	for _, kind := range kinds {
		if t.Kind == kind {
			p.walker.Next()
			return t
		}
	}

	log.Fatalf("Unexpected %d", t.Kind)
	return lex.Token{}
}

// Returns true if got token, and false if not
// func (p *Parser) Optional(kind lex.TokenKind) (lex.Token, bool) {
// 	t := p.walker.Get()
// 	if t.Kind == kind {
// 		p.walker.Next()
// 		return t, true
// 	}
// 	return t, false
// }

// Checks if the current tokens is any of the provided without moving the walker
func (p *Parser) IsCurrent(kinds ...lex.TokenKind) bool {
	t := p.walker.Get()
	for _, kind := range kinds {
		if t.Kind == kind {
			return true
		}
	}

	return false
}

// Returns true if got token, and false if not

func (p *Parser) Optional(kinds ...lex.TokenKind) (lex.Token, bool) {
	t := p.walker.Get()
	for _, kind := range kinds {
		if t.Kind == kind {
			p.walker.Next()
			return t, true
		}
	}

	return lex.Token{}, false
}
