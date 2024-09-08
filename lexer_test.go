package main

import (
	"bufio"
	"strings"
	"testing"

	"github.com/The-Ocean-Man/hightide-c/lexer"
)

func generateSingleToken(str string) lexer.Token {
	cs := lexer.NewCharStream(bufio.NewReader(strings.NewReader(str)))
	l := lexer.NewLexer(&cs)
	p := l.Parse()

	// These panics should be unreachable
	if len(p) != 1 {
		panic("Fuckup1")
	}
	if len(p[0].Children) != 0 {
		panic("Fuckup2")
	}
	if len(p[0].Content) != 1 {
		panic("Fuckup3")
	}

	return p[0].Content[0]
}

func Test_Tokens1(t *testing.T) {
	assert(t, generateSingleToken("mod").Kind == lexer.MOD, "mod failed")
	assert(t, generateSingleToken("use").Kind == lexer.USE, "use failed")
	assert(t, generateSingleToken("alias").Kind == lexer.ALIAS, "alias failed")
	assert(t, generateSingleToken("if").Kind == lexer.IF, "if failed")
	assert(t, generateSingleToken("else").Kind == lexer.ELSE, "else failed")
	assert(t, generateSingleToken("return").Kind == lexer.RETURN, "return failed")
	assert(t, generateSingleToken("struct").Kind == lexer.STRUCT, "struct failed")
	assert(t, generateSingleToken("enum").Kind == lexer.ENUM, "enum failed")
	assert(t, generateSingleToken("var").Kind == lexer.VAR, "var failed")
	assert(t, generateSingleToken("const").Kind == lexer.CONST, "const failed")
	assert(t, generateSingleToken("rdo").Kind == lexer.RDO, "rdo failed")
	assert(t, generateSingleToken("ref").Kind == lexer.REF, "ref failed")
	assert(t, generateSingleToken("attr").Kind == lexer.ATTR, "attr failed")
	assert(t, generateSingleToken("async").Kind == lexer.ASYNC, "async failed")
	assert(t, generateSingleToken("do").Kind == lexer.DO, "do failed")
	assert(t, generateSingleToken("for").Kind == lexer.FOR, "for failed")

	hello := generateSingleToken("hello_World123")
	assert(t, hello.Kind == lexer.NAME && hello.Data.(string) == "hello_World123", "ident failed")
	str := generateSingleToken("\"this is so much fun 123 _ @£$€@\"")
	assert(t, str.Kind == lexer.STRING && str.Data.(string) == "this is so much fun 123 _ @£$€@", "string failed")

	// Parsing numbers are already tested in arithmetic tests
}
