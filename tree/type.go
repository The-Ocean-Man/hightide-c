package tree

type Type interface{}

type StructType struct {
	Fields []*StructField
}

type StructField struct {
	Name         string
	Type         Type
	DefaultValue Node // Optional
}

// type TypeKind int

// const (
// 	TKType TypeKind = iota
// 	TKPtr
// 	TKRef
// 	TKRdo
// 	TKConst
// 	TKOwnPtr
// )

func BuiltinTy(ty BuiltinType) Type {
	return &ty
}

// Just support builtin types for now
type BuiltinType int

const (
	BITYVoid BuiltinType = iota
	BITYI8
	BITYI32
	BITYF32
	BITYBool
)

var BuiltinTypeMap = map[string]BuiltinType{
	"i8":   BITYI8,
	"int":  BITYI32,
	"i32":  BITYI32,
	"f32":  BITYF32,
	"bool": BITYBool,
}

type PtrType struct {
	Inner Type
}

type RDOType struct {
	Inner Type
}

type RefType struct {
	Inner Type
}
