package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/The-Ocean-Man/hightide-c/ast"
	"github.com/The-Ocean-Man/hightide-c/lexer"
	"github.com/The-Ocean-Man/hightide-c/parser"
	"github.com/The-Ocean-Man/hightide-c/tree"
)

func main() {
	// bytes, err := os.ReadFile("./test.ht")
	file, err := os.Open("./test.ht")
	r := bufio.NewReader(file)

	if err != nil {
		log.Fatalln(err)
	}

	// str := string(bytes)
	stream := lexer.NewCharStream(r)
	l := lexer.NewLexer(&stream)

	lines := l.Parse()
	// fmt.Println(lines[0])
	// lexer.PrintTree(lines)
	p := parser.MakeParser(lines)
	prog := p.ParseProgram()
	// fmt.Println(printExpr(prog.Children[0]))

	mod := tree.MakeModule([]*ast.ProgramNode{prog})

	tree.DiscoverTypes(mod)
	tree.DiscoverFunctions(mod)
	tree.DiscoverFuncArgs(mod)
	tree.PopulateFunctions(mod)

	// fmt.Println(*mod.Functions[0].Args["argc"].(*tree.BuiltinType))
	// fmt.Println(mod.Functions[0].Body.Children[0].(*tree.VarDec).Val)
	// fmt.Println(mod.Functions[0].Body.Children[1].(*tree.VarDec).Val)
	// fmt.Println(mod.Functions[0].Body.Children[2].(*tree.VarDec).Val) // .(tree.BinaryOpValue).Right.(tree.BinaryOpValue).Right.(tree.Ident).Child.(*tree.VarDec).Name
	// fmt.Println(mod.Functions[0].Body.Children[3].(*tree.FuncCall).Function.(*tree.Ident).Child.(*tree.FuncDec).ReturnType.(*tree.BuiltinType))

	// for _, dec := range prog.Children {
	// 	if v, ok := dec.(*ast.VarDecNode); ok {
	// 		_ = v
	// 		// fmt.Println(v)
	// 	} else if b, ok := dec.(*ast.BinaryOperatorNode); ok {
	// 		// fmt.Println(b.Left.(*ast.BinaryOperatorNode).Right.(*ast.UnaryOperatorNode).Child) // b.Left.(*ast.BinaryOperatorNode).Right.(*ast.BinaryOperatorNode).Right
	// 		fmt.Println(evalExpr(b))
	// 	}
	// }
}

func evalExpr(n ast.Node) float64 {
	if b, ok := n.(*ast.BinaryOperatorNode); ok {
		switch b.GetKind() {
		case ast.NKBinaryAdd:
			return evalExpr(b.Left) + evalExpr(b.Right)
		case ast.NKBinarySub:
			return evalExpr(b.Left) - evalExpr(b.Right)
		case ast.NKBinaryMul:
			return evalExpr(b.Left) * evalExpr(b.Right)
		case ast.NKBinaryDiv:
			return evalExpr(b.Left) / evalExpr(b.Right)
		case ast.NKBinaryRem:
			panic("Remainders not supported for floats (only for ints)")
		}
	}
	if u, ok := n.(*ast.UnaryOperatorNode); ok {
		switch u.GetKind() {
		case ast.NKUnaryNegate:
			n := evalExpr(u.Child)
			if n == 0 {
				return 1
			} else {
				return 0
			}
		case ast.NKUnaryInvert:
			return -evalExpr(u.Child)
		case ast.NKUnaryRDO:
			return 100 * evalExpr(u.Child)
		}
	}
	if i, ok := n.(*ast.IntLitteralNode); ok {
		return float64(i.Value)
	}
	if f, ok := n.(*ast.FloatLitteralNode); ok {
		return f.Value
	}
	panic(fmt.Sprintln("Unsupported expr type", n.GetKind()))
}

func printExpr(e ast.Node) string {
	if bin, ok := e.(*ast.BinaryOperatorNode); ok {
		if bin.Kind == ast.NKBinaryAdd {
			return fmt.Sprintf("(%s + %s)", printExpr(bin.Left), printExpr(bin.Right))
		}
		if bin.Kind == ast.NKBinarySub {
			return fmt.Sprintf("(%s - %s)", printExpr(bin.Left), printExpr(bin.Right))
		}
		if bin.Kind == ast.NKBinaryMul {
			return fmt.Sprintf("(%s * %s)", printExpr(bin.Left), printExpr(bin.Right))
		}
	}
	if un, ok := e.(*ast.UnaryOperatorNode); ok {
		if un.Kind == ast.NKUnaryInvert {
			return fmt.Sprintf("-%s", printExpr(un.Child))
		}
	}
	if i, ok := e.(*ast.IntLitteralNode); ok {
		return fmt.Sprint(i.Value)
	}
	return "<>"
}
