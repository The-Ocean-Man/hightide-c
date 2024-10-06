package tree

import (
	"fmt"
	"reflect"

	"github.com/The-Ocean-Man/hightide-c/ast"
	"github.com/golang-collections/collections/stack"
)

type Module struct {
	GlobalVars []Node
	Functions  []*FuncDec
	Types      []Type
	//Imports    []whatever // TODO

	Files []*ast.ProgramNode
}

func (mod *Module) ResolveType(n ast.Node) Type {
	if un, ok := n.(*ast.UnaryOperatorNode); ok {
		if un.GetKind() == ast.NKUnaryPtrTo {
			return &PtrType{Inner: mod.ResolveType(un.Child)}
		}
		if un.GetKind() == ast.NKUnaryRefTo {
			return &PtrType{Inner: &RDOType{Inner: mod.ResolveType(un.Child)}}
		}
		if un.GetKind() == ast.NKUnaryREF {
			return &RefType{Inner: mod.ResolveType(un.Child)}
		}
	}

	return mod.findType(n)
}

func (mod *Module) findType(n ast.Node) Type {
	if n.GetKind() != ast.NKIdent {
		panic("TODO: add anonymous types")
	}
	tyName := n.(*ast.IdentNode).Name

	// First check against built in types
	for name, bty := range BuiltinTypeMap {
		if tyName == name {
			return BuiltinTy(bty)
		}
	}
	fmt.Println(n.(*ast.IdentNode).GetString())

	// Check against defined types
	for _, ty := range mod.Types {
		_ = ty // TODO: implement once user defined types are implemented
	}

	panic(fmt.Sprintln("Type resolution failed:", n.GetKind()))
}

func MakeModule(files []*ast.ProgramNode) *Module {
	mod := &Module{make([]Node, 0), make([]*FuncDec, 0), make([]Type, 0), files}

	return mod
}

func DiscoverTypes(mod *Module) {
	for _, file := range mod.Files {
		for _, top := range file.Children {
			// Discover user declared types once implemented
			_ = top
		}
	}
}

func PopulateTypes(mod *Module) {
	// TODO: discover the fields of structs and enums, now that we know all defined types
}

func DiscoverFunctions(mod *Module) {
	for _, file := range mod.Files {
		for _, top := range file.Children {
			if dec, ok := top.(*ast.FuncDecNode); ok {
				mod.Functions = append(mod.Functions, &FuncDec{Ast: dec, Name: dec.Proto.Name, Args: make(map[string]Type),
					ReturnType: nil, Body: nil, IsExtern: dec.IsExtern})
			}
		}
	}
}

// Boths args but also other attributes
func DiscoverFuncArgs(mod *Module) {
	for _, fn := range mod.Functions {
		fn.Name = fn.Ast.Proto.Name

		if fn.Ast.Proto.ReturnTy != nil {
			fn.ReturnType = mod.ResolveType(fn.Ast.Proto.ReturnTy)
		} else {
			fn.ReturnType = BuiltinTy(BITYVoid)
		}

		fn.IsExtern = fn.Ast.IsExtern

		for _, a := range fn.Ast.Proto.Args {
			arg := a.(*ast.VarDecNode)
			fn.Args[arg.Name] = mod.ResolveType(arg.Type)
		}
	}
}

func PopulateFunctions(mod *Module) { // Parse contents of functions TODO
	for _, f := range mod.Functions {
		ParseFunction(f, mod)
	}
}

func ParseFunction(fn *FuncDec, mod *Module) {
	a := fn.Ast
	ctx := stack.New()
	for name, ty := range fn.Args {
		vardec := &VarDec{Ast: nil, Mut: ast.Mutable, Name: name, Ty: ty, Val: nil, Phantom: true}
		ctx.Push(vardec)
	}

	if fn.IsExtern {
		return
	}

	block := ParseBlock(a.Body, ctx, mod)
	fn.Body = block
}

func ParseBlock(block *ast.BlockNode, ctx *stack.Stack, mod *Module) *Block {
	statements := []Node{}

	stackSize := ctx.Len()
	for _, stmt := range block.Children {
		if vardec, ok := stmt.(*ast.VarDecNode); ok {
			vd := &VarDec{Ast: vardec, Name: vardec.Name, Mut: vardec.Mut,
				Ty: mod.ResolveType(vardec.Type), Val: ParseValue(vardec.Value, ctx, mod), Phantom: false}

			ctx.Push(vd)
			statements = append(statements, vd)
			continue
		}
		if fncall, ok := stmt.(*ast.FuncCallNode); ok {
			fn := ParseValue(fncall.Func, ctx, mod)
			args := make([]Value, len(fncall.Args))
			for i, arg := range fncall.Args {
				args[i] = ParseValue(arg, ctx, mod)
			}
			call := &FuncCall{Ast: fncall, Function: fn, Args: args}

			ctx.Push(call)
			statements = append(statements, call)
			continue
		}
	}
	for stackSize != ctx.Len() {
		ctx.Pop()
	}
	return &Block{Ast: block, Children: statements}
}

func ParseValue(n ast.Node, ctx *stack.Stack, mod *Module) Value {
	if ident, ok := n.(*ast.IdentNode); ok {
		return &Ident{Child: resolveIdent(ident, ctx, mod), Ast: ident, Mut: ast.Mutable}
	}
	if deref, ok := n.(*ast.DerefNode); ok {
		val := ParseValue(deref.Inner, ctx, mod)
		return &UnaryOpValue{Ast: deref, Mut: val.GetMutability(), Child: val, Op: ast.NKDeref}
	}

	if bin, ok := n.(*ast.BinaryOperatorNode); ok {
		return &BinaryOpValue{
			Left:  ParseValue(bin.Left, ctx, mod),
			Right: ParseValue(bin.Right, ctx, mod),
			Ast:   bin,
			Op:    bin.Kind,
			Mut:   ast.Mutable,
		}
	}
	if un, ok := n.(*ast.UnaryOperatorNode); ok {
		return &UnaryOpValue{
			Ast:   un,
			Child: ParseValue(un, ctx, mod),
			Op:    un.Kind,
			Mut:   ast.Mutable,
		}
	}
	if num, ok := n.(*ast.IntLitteralNode); ok {
		return &NumberValue{
			Ast:   num,
			Child: num.Value,
			Mut:   ast.Constant,
		}
	}

	panic(fmt.Sprintf("Unknown value type '%s'", reflect.TypeOf(n)))
}

func resolveIdent(name *ast.IdentNode, ctx *stack.Stack, mod *Module) IdentValue {
	// search stack (oh god)
	backlog := stack.New()

	var foundIdent IdentValue = nil

	for {
		top := ctx.Peek()
		if top == nil {
			break // Look in global decs next
		}
		if vardec, ok := top.(*VarDec); ok {
			if vardec.Name == name.Name {
				foundIdent = vardec
				break
			}
		} else {
			if top.(*FuncDec).Name == name.Name {
				foundIdent = top.(Node)
				break
			}
		}
		backlog.Push(ctx.Pop())
	}
	for backlog.Len() != 0 {
		ctx.Push(backlog.Pop())
	}

	if foundIdent != nil {
		return foundIdent
	}

	// Search for functions in module
	// TODO: search in other modules
	for _, fn := range mod.Functions {
		if fn.Name == name.Name {
			return fn
		}
	}

	panic(fmt.Sprintf("Could not resolve '%s'", name.Name))
}
