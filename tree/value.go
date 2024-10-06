package tree

import "github.com/The-Ocean-Man/hightide-c/ast"

type Value interface {
	Node
	// GetValueKind() ValueType
	GetMutability() ast.Mutability
}

type UnaryOpType uint8

const (
	UOTDeref UnaryOpType = iota
	UOTRDO
	UOTCONST
	UOTREF
	UOTOwnPtr
	UOTAddrOf
	UOTNegate // -
	UOTInvert // !
)

type BinaryOpType uint8

const (
	BOTAdd BinaryOpType = iota
	BOTSub
	BOTMul
	BOTDiv
	BOTRem

	// TODO: Add all other bitwise operations
	BOTAnd // TODO
	BOTOr  // TODO

	BOTEq
	BOTNeq
	BOTGt
	BOTGeq
	BOTLt
	BOTLeq
	BOTBoolAnd
	BOTBoolOr
)

// type ValueType uint8

// const ()

var _ Value = &InvocationValue{}

type InvocationValue struct {
	Ast *ast.FuncCallNode

	Func Value
	Args []Value

	Mut ast.Mutability
}

func (i InvocationValue) GetAST() ast.Node {
	return i.Ast
}

func (InvocationValue) GetKind() PNodeKind {
	return PNKFuncCall
}

func (b InvocationValue) GetMutability() ast.Mutability {
	return b.Mut
}

type BinaryOpValue struct {
	Ast ast.Node

	Left  Value
	Right Value
	Op    ast.UnlinkedNodeKind

	Mut ast.Mutability
}

func (BinaryOpValue) GetKind() PNodeKind {
	return PNKBinaryOp
}

func (b BinaryOpValue) GetMutability() ast.Mutability {
	return b.Mut
}
func (v BinaryOpValue) GetAST() ast.Node {
	return v.Ast
}

type UnaryOpValue struct {
	Ast ast.Node

	Child Value
	Op    ast.UnlinkedNodeKind // Why not

	Mut ast.Mutability
}

func (UnaryOpValue) GetKind() PNodeKind {
	return PNKUnaryOp
}

func (b UnaryOpValue) GetMutability() ast.Mutability {
	return b.Mut
}
func (v UnaryOpValue) GetAST() ast.Node {
	return v.Ast
}

type IdentValue interface{} // Pointer to some declaration
// Either FuncDec or VarDec

type Ident struct {
	Ast ast.Node

	Child IdentValue // TODO: Make actual ident, idk might not need fix

	Mut ast.Mutability
}

func (Ident) GetKind() PNodeKind {
	return PNKIdent
}

func (b Ident) GetMutability() ast.Mutability {
	return b.Mut
}
func (v Ident) GetAST() ast.Node {
	return v.Ast
}

type NumberValue struct {
	Ast ast.Node

	Child int64

	Mut ast.Mutability
}

func (NumberValue) GetKind() PNodeKind {
	return PNKNumber
}

func (n NumberValue) GetMutability() ast.Mutability {
	return n.Mut
}

func (v NumberValue) GetAST() ast.Node {
	return v.Ast
}
