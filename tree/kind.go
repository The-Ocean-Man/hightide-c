package tree

type PNodeKind uint

const (
	PNKFuncDec PNodeKind = iota
	PNKBlock
	PNKVarDec
	PNKFuncCall
	PNKBinaryOp
	PNKUnaryOp
	PNKAssign
	PNKIdent
	PNKNumber
	PNKDecimal
	PNKString
)
