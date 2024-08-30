package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/The-Ocean-Man/hightide-c/ast"
	"github.com/The-Ocean-Man/hightide-c/lexer"
	"github.com/The-Ocean-Man/hightide-c/parser"
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

	for _, dec := range prog.Children {
		if v, ok := dec.(*ast.VarDecNode); ok {
			_ = v
			// fmt.Println(v)
		} else if b, ok := dec.(*ast.BinaryOperatorNode); ok {
			// fmt.Println(b.Left.(*ast.BinaryOperatorNode).Right.(*ast.UnaryOperatorNode).Child) // b.Left.(*ast.BinaryOperatorNode).Right.(*ast.BinaryOperatorNode).Right
			fmt.Println(evalExpr(b))
		}
	}
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
			return 10 * evalExpr(u.Child)
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
