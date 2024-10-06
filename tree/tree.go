package tree

import (
	"github.com/The-Ocean-Man/hightide-c/ast"
)

type Node interface {
	GetKind() PNodeKind
	GetAST() ast.Node
}

type FuncDec struct {
	Ast *ast.FuncDecNode

	Name       string
	Args       map[string]Type
	ReturnType Type
	Body       *Block // If nil then this is a extern func declaration and not impl
	IsExtern   bool
}

func (FuncDec) GetKind() PNodeKind {
	return PNKFuncDec
}
func (n *FuncDec) GetAST() ast.Node {
	return n.Ast
}

type Block struct {
	Ast *ast.BlockNode

	Children []Node
}

func (Block) GetKind() PNodeKind {
	return PNKFuncDec
}
func (n *Block) GetAST() ast.Node {
	return n.Ast
}

type VarDec struct {
	Ast *ast.VarDecNode

	Mut  ast.Mutability
	Name string
	Ty   Type
	Val  Value

	Phantom bool // Denotes if this is a variable which is passed as an argument or such, or declared formally
}

func (VarDec) GetKind() PNodeKind {
	return PNKVarDec
}
func (n *VarDec) GetAST() ast.Node {
	return n.Ast
}

type FuncCall struct {
	Ast *ast.FuncCallNode

	Function Value
	Args     []Value
}

func (FuncCall) GetKind() PNodeKind {
	return PNKFuncCall
}
func (n *FuncCall) GetAST() ast.Node {
	return n.Ast
}

type Assignment struct {
	Ast *ast.AssignmentNode

	Target Value
	Value  Value
}

func (Assignment) GetKind() PNodeKind {
	return PNKAssign
}
func (n *Assignment) GetAST() ast.Node {
	return n.Ast
}
