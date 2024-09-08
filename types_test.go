package main

import "testing"

func Test_Types1(t *testing.T) {
	p := generateProgramUAST("*int")
	assert(t, len(p.Children) == 1, "Ptr to int failed1")

}
