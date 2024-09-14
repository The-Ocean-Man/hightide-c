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

	for _, ln := range p.lines { // Parse Top Level
		walker := MakeWalker(ln)
		p.walker = walker
		if isVarDecToken(walker.Get().Kind) {
			topLevelStatements = append(topLevelStatements, p.ParseVarDec())
		} else if walker.Get().Kind == lex.NAME {
			topLevelStatements = append(topLevelStatements, p.ParseFuncDec())
		} else {
			panic("Unexpected token in toplevel")
		}
		// else if isExprStartToken(walker.Get().Kind) { // debug only, TODO: Remove
		// 	p.Expect(lex.RETURN)
		// 	topLevelStatements = append(topLevelStatements, p.ParseExpr())
		// }
	}

	return &ast.ProgramNode{Children: topLevelStatements}
}

// Returns true if kind == VAR, CONST or RDO
func isVarDecToken(kind lex.TokenKind) bool {
	return kind == lex.VAR || kind == lex.CONST || kind == lex.RDO
}

func (p *Parser) ParseFuncDec() *ast.FuncDecNode {
	proto := p.ParseFuncProto()
	if len(p.walker.ln.Children) == 0 {
		return &ast.FuncDecNode{Proto: proto, Body: nil}
	}

	block := &ast.BlockNode{Children: make([]ast.Node, 0)}
	for _, ln := range p.walker.ln.Children {
		p.walker = MakeWalker(ln)
		stmt := p.ParseStatement()
		block.Children = append(block.Children, stmt)
	}
	return &ast.FuncDecNode{Proto: proto, Body: block}
}

func (p *Parser) ParseFuncProto() *ast.FuncProtoNode {
	fnName := p.Expect(lex.NAME).Data.(string)

	// Parse args
	args := make([]ast.Node, 0)
	p.Expect(lex.LPAREN) // generics come later
	for {
		if p.IsCurrent(lex.EOF, lex.EOL) {
			panic("Expected end of function proto")
		}
		if p.OptionalNoVal(lex.RPAREN) {
			break
		}
		argName := p.Expect(lex.NAME).Data.(string)
		argTy := p.ParseExpr()
		args = append(args, &ast.VarDecNode{Name: argName, Type: argTy, Mut: ast.Mutable, Value: nil}) // maybe change mutable to readonly

		if p.OptionalNoVal(lex.COMMA) {
			continue
		}
	}

	// Parse return
	if p.IsCurrent(lex.EOL, lex.EOF) {
		return &ast.FuncProtoNode{Name: fnName, Args: args, ReturnTy: nil}
	}

	returnTy := p.ParseExpr()
	return &ast.FuncProtoNode{Name: fnName, Args: args, ReturnTy: returnTy}
}

// Anything inside of a function block
func (p *Parser) ParseStatement() ast.Node {
	c := p.walker.Get().Kind
	if isExprStartToken(c) {
		return p.ParseExpr()
	} else if isVarDecToken(c) {
		return p.ParseVarDec()
	}

	panic(fmt.Sprintln("Unexpected token in statement:", c, p.walker.Get().Data))
}

func (p *Parser) ParseVarDec() *ast.VarDecNode {
	decTokKind := p.ExpectAnyOf(lex.VAR, lex.CONST, lex.RDO).Kind
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
	bit := p.ParseExprTerminal()

	if _, isNode := bit.(ast.Node); !isNode {
		return bit
	}
	for {
		if p.OptionalNoVal(lex.LPAREN) { // func invocation
			if callee, ok := bit.(ast.Node); ok {
				args := p.ParseArgsListWithTerm(lex.RPAREN)
				bit = &ast.FuncCallNode{Func: callee, Args: args}
			}
		} else if p.OptionalNoVal(lex.LSQUARE) { // indexing
			if outer, ok := bit.(ast.Node); ok {
				args := p.ParseArgsListWithTerm(lex.RSQUARE)
				bit = &ast.IndexNode{Outer: outer, Inner: args}
			} else {
				panic("Expected expression as index")
			}
		} else if tok, found := p.Optional(lex.DOT, lex.COLONCOLON); found {
			isDot := tok.Kind == lex.DOT

			if p.OptionalNoVal(lex.STAR) {
				if !isDot {
					panic("Cannot use foo::* syntax")
				}
				bit = &ast.DerefNode{Inner: bit.(ast.Node)} // bit is verified as ast.Node
				continue
			}

			inner := p.ParseIdent()
			bit = &ast.PropertyIndexNode{Outer: bit.(ast.Node), Inner: inner, IsDot: isDot}
		} else {
			break
		}
	}

	// account for invocations and indexing
	return bit
}

// Long name, parser the arguments of a function call
func (p *Parser) ParseArgsListWithTerm(terminator lex.TokenKind) []ast.Node {
	args := make([]ast.Node, 0)

	for {
		if _, ok := p.Optional(terminator); ok {
			break
		}

		arg := p.ParseExpr()
		args = append(args, arg)

		tok := p.ExpectAnyOf(terminator, lex.COMMA)
		if tok.Kind == terminator {
			break
		}
	}
	return args
}

func (p *Parser) ParseExprTerminal() ExprBit {
	// current := p.walker.Get()
	if p.IsCurrent(lex.EOL, lex.EOF, lex.RPAREN, lex.RSQUARE, lex.COMMA) { // end of expr
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

func (p *Parser) ParseIdent() ast.Node {
	name := p.Expect(lex.NAME)
	// top := &ast.IdentNode{}
	// top.Name = first.Data.(string)
	// var node ast.InnerSettable = top

	// for {
	// 	if p.OptionalNoVal(lex.DOT) {
	// 		node.UsedDot = true
	// 	} else if p.OptionalNoVal(lex.COLONCOLON) {
	// 		node.UsedDot = false
	// 	} else {
	// 		break
	// 	}

	// 	name := p.Expect(lex.NAME)
	// 	lower := &ast.IdentNode{}
	// 	lower.Name = name.Data.(string)
	// 	node.SetInner(lower)

	// 	node = lower
	// }

	return &ast.IdentNode{Name: name.Data.(string)}
}

func (p *Parser) Expect(kind lex.TokenKind) lex.Token {
	t := p.walker.Get()
	if t.Kind == kind {
		p.walker.Next()
		return t
	}
	panic(fmt.Sprintf("Expected %d but got %d instead", kind, t.Kind))
}

func (p *Parser) ExpectAnyOf(kinds ...lex.TokenKind) lex.Token {
	t := p.walker.Get()
	for _, kind := range kinds {
		if t.Kind == kind {
			p.walker.Next()
			return t
		}
	}

	panic(fmt.Sprintf("Unexpected %d", t.Kind))

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

func (p *Parser) OptionalNoVal(kinds ...lex.TokenKind) bool {
	_, found := p.Optional(kinds...)
	return found
}

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
