package parser

// Expression parser
import (
	"fmt"
	"log"

	"github.com/FridaFino/goalgorithms/structure/dynamicarray"
	"github.com/The-Ocean-Man/hightide-c/ast"
	"github.com/The-Ocean-Man/hightide-c/lexer"
)

const ( // Operator precedence
	PREC_NEGATE = 10
	PREC_INVOKE // calling function
	PREC_INDEX

	PREC_PLUS = 20
	PREC_MINUS

	PREC_TIMES = 30
	PREC_DIV
	PREC_MODULO

	PREC_REF = 40
	PREC_OWNPTR
	PREC_RDO
	PREC_CONST
)

type ExprBit any // lexer.TokenKind || ast.Node(another expr basically). since Node is another 'any'm if ExprBit is not NodeKind, it is Node

//#region
// func ParseExpression(bits []ExprBit) ast.Node {
// 	fmt.Println(bits)
// 	b2 := make([]ExprBit, 0)
// 	for idx, bit := range bits {
// 		if kind, ok := bit.(lexer.TokenKind); ok {
// 			if kind == lexer.PLUS {
// 				left := bits[idx-1].(ast.Node)
// 				right := bits[idx+1].(ast.Node)
// 				b2 = append(b2, ast.BinaryOperatorNode{Kind: ast.NKBinaryAdd, Left: left, Right: right})
// 				continue
// 			}
// 			b2 = append(b2, bit) // if we dont change it we ignore it
// 		}

// 	}

// 	b3 := make([]ExprBit, 0)
// 	for idx, bit := range b2 {
// 		if kind, ok := bit.(lexer.TokenKind); ok {
// 			if kind == lexer.STAR {
// 				left := bits[idx-1].(ast.Node)
// 				right := bits[idx+1].(ast.Node)
// 				b3 = append(b3, ast.BinaryOperatorNode{Kind: ast.NKBinaryMul, Left: left, Right: right})
// 				continue
// 			}
// 			b3 = append(b3, bit) // if we dont change it we ignore it
// 		}
// 	}

// 	if len(b3) != 1 {
// 		log.Fatalf("Houston we have a problem %d", len(b3))
// 	}
// 	fmt.Println(b3)

// 	return b3[0].(ast.Node)
// }
//#endregion

func ParseExpression(bits []ExprBit) ast.Node {
	da := dynamicarray.DynamicArray{}

	for _, bit := range bits {
		da.Add(bit)
	}

	// do in reverse order of precedence, but unarys first

	// negate
	iterBitsUnary(&da, func(left, right ExprBit) ExprBit {
		if tok, ok := left.(lexer.TokenKind); ok && tok == lexer.BANG { // bang is tmp, use minus
			if n, ok := right.(ast.Node); ok {
				fmt.Println(n)
				return &ast.UnaryOperatorNode{Child: n, Kind: ast.NKUnaryNegate}
			} else {
				log.Fatalf("Expected value after negation but got: %d\n", tok)
			}
		}
		return nil
	})

	// times, div, mod
	iterBitsBinary(&da, func(left, middle, right ExprBit) ExprBit {
		if tok, ok := middle.(lexer.TokenKind); ok && tok == lexer.STAR || tok == lexer.SLASH || tok == lexer.PERCENT {
			if _, ok := left.(ast.Node); !ok {
				log.Fatalln("Expected value before operator but got", left)
			}
			if _, ok := right.(ast.Node); !ok {
				log.Fatalln("Expected value after operator but got", right)
			}
			var kind ast.UnlinkedNodeKind
			if tok == lexer.STAR {
				kind = ast.NKBinaryMul
			} else if tok == lexer.SLASH {
				kind = ast.NKBinaryDiv
			} else {
				kind = ast.NKBinaryRem
			}
			return &ast.BinaryOperatorNode{Kind: kind, Left: left.(ast.Node), Right: right.(ast.Node)}
		}
		return nil
	})

	// plus minus
	iterBitsBinary(&da, func(left, middle, right ExprBit) ExprBit {
		if tok, ok := middle.(lexer.TokenKind); ok && tok == lexer.PLUS || tok == lexer.MINUS {

			if _, ok := left.(ast.Node); !ok {
				log.Fatalln("Expected value before operator but got", left)
			}
			if _, ok := right.(ast.Node); !ok {
				log.Fatalln("Expected value after operator but got", right)
			}
			var kind ast.UnlinkedNodeKind
			if tok == lexer.PLUS {
				kind = ast.NKBinaryAdd
			} else {
				kind = ast.NKBinarySub
			}
			return &ast.BinaryOperatorNode{Kind: kind, Left: left.(ast.Node), Right: right.(ast.Node)}
		}
		fmt.Println("nay")
		return nil
	})

	if da.Size != 1 {
		log.Fatalf("Expected expr list size to be one but was %d\n", da.Size)
	}

	s, ok := da.GetData()[0].(ast.Node)
	if !ok {
		log.Fatalln("Expected expr list to be of type ast.Node")
	}

	return s
}

// Return nil if nothing was detected, else return a value which will replace all its components
func iterBitsUnary(da *dynamicarray.DynamicArray, action func(left ExprBit, right ExprBit) ExprBit) {
	var idx int = 0
	for {
		if da.Size < 2 || idx+1 >= da.Size {
			break
		}

		val := action(noerr(da.Get(idx)), noerr(da.Get(idx+1)))
		if val == nil {
			idx++
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
func noerr(a any, _ error) any {
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
