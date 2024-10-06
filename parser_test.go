package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/The-Ocean-Man/hightide-c/ast"
	"github.com/The-Ocean-Man/hightide-c/lexer"
	"github.com/The-Ocean-Man/hightide-c/parser"
)

func generateProgramUAST(src string) *ast.ProgramNode {
	r := bufio.NewReader(strings.NewReader(src))
	stream := lexer.NewCharStream(r)
	l := lexer.NewLexer(&stream)

	lines := l.Parse()
	p := parser.MakeParser(lines)
	prog := p.ParseProgram()
	return prog
}

func Test_ArithNumerics1(t *testing.T) {
	return
	prog := generateProgramUAST(`(6 + 8) * -(-3+5)/2`)

	if evalExpr(prog.Children[0]) != ((6 + 8) * -(-3 + 5) / 2) { // - 14
		t.Fatal(evalExpr(prog.Children[0]))
	}
}

func Test_ArithNumerics2(t *testing.T) {
	return

	prog := generateProgramUAST(`123 * -(-192 + -98.432) * 8 /3.2`)
	assert(t, len(prog.Children) == 1, "Expr len should be one")
	expr := prog.Children[0]
	assert(t, evalExpr(expr) == (123*-(-192+-98.432)*8/3.2), "Expr results does not equal")
}

func Test_Expression1(t *testing.T) {
	return

	prog := generateProgramUAST(`hello.funny.*`)

	if len(prog.Children) != 1 {
		t.Fatal("should have generated only one line")
	}

	if deref, ok := prog.Children[0].(*ast.DerefNode); ok {
		if prop, ok := deref.Inner.(*ast.PropertyIndexNode); ok {
			if hello, ok := prop.Outer.(*ast.IdentNode); ok {
				if hello.Name != "hello" {
					t.Fatal("'hello' ident's name was not \"hello\":", hello.Name)
				}
			} else {
				t.Fatal("Hello was not ident")
			}
			if funny, ok := prop.Inner.(*ast.IdentNode); ok {
				if funny.Name != "funny" {
					t.Fatal("'funny' ident's name was not \"funny\":", funny.Name)
				}
			} else {
				t.Fatal("Hello was not ident")
			}
		} else {
			t.Fatal("prop was not prop index type")
		}
	} else {
		t.Fatal("Token was not of type deref")
	}
}

func Test_Expression2(t *testing.T) {
	return

	prog := generateProgramUAST(`silly::massive.what(123).*.omega`)

	if len(prog.Children) != 1 {
		t.Fatal("should have generated only one line")
	}

	topuncasted := prog.Children[0]
	top := tyAssert[*ast.PropertyIndexNode](t, topuncasted)

	omega := tyAssert[*ast.IdentNode](t, top.Inner)
	assert(t, omega.Name == "omega", "Expected ident omega")

	deref := tyAssert[*ast.DerefNode](t, top.Outer)

	fn := tyAssert[*ast.FuncCallNode](t, deref.Inner)
	assert(t, len(fn.Args) == 1, "Len of args should be one")

	arg := tyAssert[*ast.IntLitteralNode](t, fn.Args[0])
	assert(t, arg.Value == 123, "Func arg was not equal to '123'")

	what := tyAssert[*ast.PropertyIndexNode](t, fn.Func)
	assert(t, what.IsDot, "What is not a dot")

	whatIdent := tyAssert[*ast.IdentNode](t, what.Inner)
	assert(t, whatIdent.Name == "what", "'what' != 'what'")

	last := tyAssert[*ast.PropertyIndexNode](t, what.Outer)
	assert(t, !last.IsDot, "Apparently '::' is the same as '.'")

	massive := tyAssert[*ast.IdentNode](t, last.Inner)
	assert(t, massive.Name == "massive", "massive name assert")

	silly := tyAssert[*ast.IdentNode](t, last.Outer)
	assert(t, silly.Name == "silly", "silly name assert")
}

func Test_Assigmnents(t *testing.T) {
	prog := generateProgramUAST(`
main()
	a = 123
	b = a
	(*b) = a + 321
	`)

	assert(t, len(prog.Children) == 1, "Parsing assigmnent func failed")
	fn := tyAssert[*ast.FuncDecNode](t, prog.Children[0])
	assert(t, len(fn.Body.Children) == 3, "Parsing assigmnent func body failed")

	a1 := tyAssert[*ast.AssignmentNode](t, fn.Body.Children[0])
	assert(t, a1.Target.(*ast.IdentNode).Name == "a", "Simple assigmnent target parsing failed")
	assert(t, a1.Value.(*ast.IntLitteralNode).Value == 123, "Simple assigmnent value parsing failed")

	a3 := tyAssert[*ast.AssignmentNode](t, fn.Body.Children[2])
	target := tyAssert[*ast.UnaryOperatorNode](t, a3.Target)
	value := tyAssert[*ast.BinaryOperatorNode](t, a3.Value)

	assert(t, target.GetKind() == ast.NKUnaryPtrTo, "Complex assigmnent target parsing failed")
	assert(t, value.GetKind() == ast.NKBinaryAdd, "Complex assigmnent value parsing failed")
}

func tyAssert[T ast.Node](t *testing.T, n ast.Node) T {
	if node, ok := n.(T); ok {
		return node
	}
	var tmp T
	t.Fatalf("Expected node to be of type %s but got %s instead.", reflect.TypeOf(tmp).String(), reflect.TypeOf(n).String())
	return tmp
}

func assert(t *testing.T, b bool, msg string) {
	if !b {
		t.Fatal(msg)
	}
}
