package parser

// Expression parser
import (
	"fmt"
	"log"
	"reflect"

	"github.com/FridaFino/goalgorithms/structure/dynamicarray"
	"github.com/The-Ocean-Man/hightide-c/ast"
	"github.com/The-Ocean-Man/hightide-c/lexer"
)

// const ( // Operator precedence
// 	PREC_NEGATE = 10
// 	PREC_INVOKE // calling function
// 	PREC_INDEX

// 	PREC_PLUS = 20
// 	PREC_MINUS

// 	PREC_TIMES = 30
// 	PREC_DIV
// 	PREC_MODULO

// 	PREC_REF = 40
// 	PREC_OWNPTR
// 	PREC_RDO
// 	PREC_CONST
// )

// Temp: this is just to test expressions, they are otherwise not allowed in top scope, as soon as funcs work, i fix this TODO
// func isExprStartToken(kind lexer.TokenKind) bool {
// 	return kind == lexer.RETURN
// }

func isExprStartToken(kind lexer.TokenKind) bool {
	return kind == lexer.NUMBER || kind == lexer.DECIMAL || kind == lexer.LPAREN || kind == lexer.MINUS || kind == lexer.BANG || // kind == lexer.
		kind == lexer.STAR || kind == lexer.AMPERSAND || kind == lexer.NAME || kind == lexer.STRING
}

type ExprBit any // lexer.TokenKind || ast.Node(another expr basically). since Node is another 'any' if ExprBit is not NodeKind, it is Node

func doUnary(left, right ExprBit, tokenKind lexer.TokenKind, nodeKind ast.UnlinkedNodeKind) (ExprBit, bool) {
	if tok, ok := left.(lexer.TokenKind); ok && tok == tokenKind { // bang is tmp, use minus
		if n, ok := right.(ast.Node); ok {
			return &ast.UnaryOperatorNode{Child: n, Kind: nodeKind}, true
		} else {
			log.Fatalf("Expected expr after unary '%c' but got: %d\n", tokenKind, tok)
		}
	}
	return nil, false
}

func doBinary(left, middle, right ExprBit, ops map[lexer.TokenKind]ast.UnlinkedNodeKind) (ExprBit, bool) {
	if tok, ok := middle.(lexer.TokenKind); ok {
		wasMatch := false
		for op := range ops {
			if op == tok {
				wasMatch = true
				break
			}
		}
		if !wasMatch {
			return nil, true
		}

		if _, ok := left.(ast.Node); !ok {
			panic(fmt.Sprintln("Expected value before operator but got", left))
		}
		if _, ok := right.(ast.Node); !ok {

			panic(fmt.Sprintln("Expected value before operator but got", right))
		}
		var kind = ops[tok]
		return &ast.BinaryOperatorNode{Kind: kind, Left: left.(ast.Node), Right: right.(ast.Node)}, true
	}
	return nil, false
}

func ParseExpression(bits []ExprBit) ast.Node {
	da := dynamicarray.DynamicArray{}

	for _, bit := range bits {
		da.Add(bit)
	}

	// do in reverse order of precedence, but unarys first

	// negate
	iterBitsUnary(&da, func(left, right ExprBit) ExprBit {
		if bit, ok := doUnary(left, right, lexer.BANG, ast.NKUnaryNegate); ok {
			return bit
		}
		if bit, ok := doUnary(left, right, lexer.MINUS, ast.NKUnaryInvert); ok {
			return bit
		}
		if bit, ok := doUnary(left, right, lexer.STAR, ast.NKUnaryPtrTo); ok {
			return bit
		}
		if bit, ok := doUnary(left, right, lexer.AMPERSAND, ast.NKUnaryRefTo); ok {
			return bit
		}
		if bit, ok := doUnary(left, right, lexer.REF, ast.NKUnaryREF); ok {
			return bit
		}
		if bit, ok := doUnary(left, right, lexer.CONST, ast.NKUnaryCONST); ok {
			return bit
		}
		if bit, ok := doUnary(left, right, lexer.RDO, ast.NKUnaryRDO); ok {
			return bit
		}
		if bit, ok := doUnary(left, right, lexer.DOLLAR, ast.NKUnaryOwnPtr); ok {
			return bit
		}
		return nil
	})

	// times, div, mod
	iterBitsBinary(&da, func(left, middle, right ExprBit) ExprBit {
		if bit, ok := doBinary(left, middle, right, map[lexer.TokenKind]ast.UnlinkedNodeKind{
			lexer.STAR:    ast.NKBinaryMul,
			lexer.SLASH:   ast.NKBinaryDiv,
			lexer.PERCENT: ast.NKBinaryRem,
		}); ok {
			return bit
		}
		return nil
	})

	// plus minus
	iterBitsBinary(&da, func(left, middle, right ExprBit) ExprBit {
		if bit, ok := doBinary(left, middle, right, map[lexer.TokenKind]ast.UnlinkedNodeKind{
			lexer.PLUS:  ast.NKBinaryAdd,
			lexer.MINUS: ast.NKBinarySub,
		}); ok {
			return bit
		}
		return nil
	})

	if da.Size != 1 {
		fmt.Println(reflect.TypeOf(da.GetData()[0]))
		log.Fatalf("Expected expr list size to be one but was %d\n", da.Size)
	}

	s, ok := da.GetData()[0].(ast.Node)
	if !ok {
		log.Fatalln("Expected expr list to be of type ast.Node")
	}

	return s
}

func bitIsTokenKindOrNil(b ExprBit) (ret bool) {
	if b == nil {
		return true
	}
	_, ret = b.(lexer.TokenKind)
	return
}

// Return nil if nothing was detected, else return a value which will replace all its components.
// Unary iteration is dont in reverse which allows things like **int and such
func iterBitsUnary(da *dynamicarray.DynamicArray, action func(left, right ExprBit) ExprBit) { // prev is the token before left, used in cases like a - -b
	var idx int = da.Size - 2

	for {
		if da.Size < 2 || idx < 0 {
			break
		}
		if !bitIsTokenKindOrNil(noerr(da.Get(idx - 1))) { // ensures that tokens like minus are only parsed as unary when they should be
			idx--
			continue
		}
		// fmt.Println(noerr(da.Get(idx)), noerr(da.Get(idx+1)))
		val := action(noerr(da.Get(idx)), noerr(da.Get(idx+1)))
		if val == nil {
			idx--
			continue
		} // else

		// remove both bits
		da.Remove(idx)
		da.Remove(idx)

		// replace bits
		// da.Put(idx, val) // dont increment idx
		pushElem(da, idx, val)
	}
}

// Return nil if nothing was detected, else return a value which will replace all its components
func iterBitsBinary(da *dynamicarray.DynamicArray, action func(left, middle, right ExprBit) ExprBit) {
	var idx int = 0
	for {
		if da.Size < 3 || idx+2 >= da.Size {
			break
		}

		val := action(noerr(da.Get(idx)), noerr(da.Get(idx+1)), noerr(da.Get(idx+2)))
		if val == nil {
			idx++
			continue
		} // else

		// remove component bits
		da.Remove(idx)
		da.Remove(idx)
		da.Remove(idx)

		// replace bits
		// da.Put(idx, val) // dont increment idx
		pushElem(da, idx, val)
	}
}
func noerr(a any, err error) any {
	if err != nil {
		return nil
	}
	return a
}

func pushElem(da *dynamicarray.DynamicArray, idx int, elem any) error {
	da.Add(nil)
	for i := da.Size - 1; i > idx; i-- {
		g, err := da.Get(i - 1)
		if err != nil {
			return err
		}
		err = da.Put(i, g)
		if err != nil {
			return err
		}
	}
	err := da.Put(idx, elem)
	return err
}
