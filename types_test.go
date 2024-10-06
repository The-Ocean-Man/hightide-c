package main

import (
	"testing"

	"github.com/The-Ocean-Man/hightide-c/ast"
)

func Test_Types1(t *testing.T) {
	p := generateProgramUAST("var variable_ *int = 123")

	assert(t, len(p.Children) == 1, "Ptr to int failed1")

	vardec := tyAssert[*ast.VarDecNode](t, p.Children[0])
	assert(t, vardec.Mut == ast.Mutable, "Variable mutability not parsing correctly")
	assert(t, vardec.Name == "variable_", "Variable name not parsing correctly")
	ptr := tyAssert[*ast.UnaryOperatorNode](t, vardec.Type)
	assert(t, ptr.Kind == ast.NKUnaryPtrTo, "Variable type not parsing correctly 1")
	Int := tyAssert[*ast.IdentNode](t, ptr.Child)
	assert(t, Int.Name == "int", "Variable type not parsing correctly 2")
}

// func Test_Types2(t *testing.T) {
// 	generateProgramUAST("*[N]&int")
// }
