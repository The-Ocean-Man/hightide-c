package main

import (
	"testing"

	"github.com/The-Ocean-Man/hightide-c/ast"
)

func Test_FuncProto1(t *testing.T) {
	prog := generateProgramUAST(`main(argc int, argv **char) int`)
	fnDec := tyAssert[*ast.FuncDecNode](t, prog.Children[0])
	assert(t, fnDec.Body == nil, "Func body was not nil when expected")
	proto := fnDec.Proto
	assert(t, proto.Name == "main", "Func name was not 'main'")
	assert(t, proto.ReturnTy.GetKind() == ast.NKIdent && proto.ReturnTy.(*ast.IdentNode).Name == "int", "Return type not parsed as expected")

	assert(t, len(proto.Args) == 2, "Args not parsed correctly")

	arg1 := tyAssert[*ast.VarDecNode](t, proto.Args[0])
	arg2 := tyAssert[*ast.VarDecNode](t, proto.Args[1])

	assert(t, arg1.Name == "argc", "Arg names not parsed correctly 1")
	assert(t, arg2.Name == "argv", "Arg names not parsed correctly 2")

	arg1Ty := tyAssert[*ast.IdentNode](t, arg1.Type)
	assert(t, arg1Ty.Name == "int", "Argc type was not int")

	arg2Ty1 := tyAssert[*ast.UnaryOperatorNode](t, arg2.Type)
	assert(t, arg2Ty1.Kind == ast.NKUnaryPtrTo, "Arg type parsed incorrectly 1")
	arg2Ty2 := tyAssert[*ast.UnaryOperatorNode](t, arg2Ty1.Child)
	assert(t, arg2Ty2.Kind == ast.NKUnaryPtrTo, "Arg type parsed incorrectly 2")

	char := tyAssert[*ast.IdentNode](t, arg2Ty2.Child)
	assert(t, char.Name == "char", "char was not char")
}

func Test_Func1(t *testing.T) {
	prog := generateProgramUAST(`
main(argc *int, argv*char) int
	var funnyVar int = 0
`)
	_ = prog

	assert(t, prog.Children[0].(*ast.FuncDecNode).Body.Children[0].(*ast.VarDecNode).Name == "funnyVar", "It dont work")
}

func Test_FuncCall1(t *testing.T) {
	prog := generateProgramUAST(`
main() 
	funny(213, 321, "Hello world", true)
`)

	assert(t, len(prog.Children) == 1, "Expected only one item in program")

	fnDec := tyAssert[*ast.FuncDecNode](t, prog.Children[0])

	proto := tyAssert[*ast.FuncProtoNode](t, fnDec.Proto)
	assert(t, proto.Name == "main", "Function name parsing incorrectly")
	assert(t, len(proto.Args) == 0, "Function args parsing incorrectly")
	assert(t, proto.ReturnTy == nil, "Function return type parsing incorrectly")

	block := fnDec.Body

	assert(t, len(block.Children) == 1, "Function body not parsing correctly")
	fnCall := tyAssert[*ast.FuncCallNode](t, block.Children[0])
	assert(t, len(fnCall.Args) == 4, "Funccall args not parsing correctly")
	assert(t, fnCall.Args[2].(*ast.StringLitteralNode).Value == "Hello world", "String in functioncall not parsing correctly")
}
