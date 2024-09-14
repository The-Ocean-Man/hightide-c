package main

import "testing"

func Test_Types1(t *testing.T) {
	p := generateProgramUAST("var i *int")
	assert(t, len(p.Children) == 1, "Ptr to int failed1")
}

// func Test_Types2(t *testing.T) {
// 	generateProgramUAST("*[N]&int")
// }
