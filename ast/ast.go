// Unlinked AST
package ast

import (
	"fmt"
)

type UnlinkedNodeKind int

const (
	NKProgram UnlinkedNodeKind = iota
	NKModule
	NKVarDec
	NKFuncCall
	NKFuncDec
	NKIndex
	NKIdent
	NKBinaryAdd
	NKBinarySub
	NKBinaryMul
	NKBinaryDiv
	NKBinaryRem
	NKUnaryNegate
	NKUnaryOwnPtr
	NKUnaryREF
	NKUnaryRDO
	NKUnaryCONST
	NKString
	NKInt
	NKFloat
)

type Node interface {
	GetKind() UnlinkedNodeKind
}

type ProgramNode struct { // One per file
	Children []Node
}

func (n ProgramNode) GetKind() UnlinkedNodeKind {
	return NKProgram
}

type FuncDecNode struct {
	Name     string
	Args     []Node
	ReturnTy Node       //optional
	Body     *BlockNode // if nil then its extern function
}

func (n FuncDecNode) GetKind() UnlinkedNodeKind {
	return NKFuncDec
}

type BlockNode struct {
	Children []Node
}

type StructNode struct {
	Name string
}

type Mutability uint8

const (
	Mutable  Mutability = iota
	ReadOnly            // RDO
	Constant            // static immutable
)

type VarDecNode struct {
	Name  string
	Mut   Mutability // This is stored in the type of the var
	Type  Node       // Optional, infer
	Value Node       // Optional
}

func (n VarDecNode) GetKind() UnlinkedNodeKind {
	return NKVarDec
}

type FuncCallNode struct {
	Func Node
	Args []Node
}

func (n FuncCallNode) GetKind() UnlinkedNodeKind {
	return NKFuncCall
}

type IndexNode struct {
	Outer Node
	Inner Node
}

func (n IndexNode) GetKind() UnlinkedNodeKind {
	return NKIndex
}

type BinaryOperatorNode struct {
	Kind  UnlinkedNodeKind // add, sub, mul, div, rem
	Left  Node
	Right Node
}

func (n BinaryOperatorNode) GetKind() UnlinkedNodeKind {
	return n.Kind
}

type UnaryOperatorNode struct {
	Kind  UnlinkedNodeKind // negate, ref, ownptr($), rdo, const
	Child Node
}

func (n UnaryOperatorNode) GetKind() UnlinkedNodeKind {
	return n.Kind
}

type IdentNode struct {
	Name    string
	Child   *IdentNode // Optional
	UsedDot bool       // true = . || false = ::
}

func (n IdentNode) GetString() string {
	var con string
	if n.UsedDot {
		con = "."
	} else {
		con = "::"
	}
	if n.Child == nil {
		return n.Name
	} else {
		return fmt.Sprintf("%s%s%s", n.Name, con, n.Child.GetString())
	}
}

func (n IdentNode) GetKind() UnlinkedNodeKind {
	return NKIdent
}

type ModuleNode struct {
	Name *IdentNode
}

func (n ModuleNode) GetKind() UnlinkedNodeKind {
	return NKModule
}

type StringLitteralNode struct {
	Value string
}

func (n StringLitteralNode) GetKind() UnlinkedNodeKind {
	return NKString
}

type IntLitteralNode struct {
	Value int64
}

func (n IntLitteralNode) GetKind() UnlinkedNodeKind {
	return NKInt
}

type FloatLitteralNode struct {
	Value float64
}

func (n FloatLitteralNode) GetKind() UnlinkedNodeKind {
	return NKFloat
}
