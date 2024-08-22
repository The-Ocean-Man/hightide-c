package lexer

import "fmt"

type Line struct {
	Content  []Token
	Children []*Line
}

type Token struct {
	Kind TokenKind
	Data any // string/number/whatever
}

// func (t Token) GetKind() TokenKind {
// 	return t.Kind
// }

func (t Token) String() (s string, b bool) {
	switch data := t.Data.(type) {
	case string:
		return data, true
	case int64:
		return fmt.Sprint(data), true
	case float64:
		return fmt.Sprint(data), true
	}

	s = fmt.Sprint(t.Kind)
	b = len(s) != 0

	return
}
func (t Token) Number() (i int64, b bool) {
	i, b = t.Data.(int64)
	return
}

func (t Token) Decimal() (f float64, b bool) {
	f, b = t.Data.(float64)
	return
}
