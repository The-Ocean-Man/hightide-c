package lexer

import (
	"fmt"
	"strings"
)

func PrintTree(lines []*Line) {
	printTreeImpl(lines, 0)
}

func printTreeImpl(lines []*Line, depth int) {
	put := func(l *Line) {
		fmt.Print(strings.Repeat(" ", depth*4))
		for _, tok := range l.Content {
			s, ok := tok.String()
			if !ok {
				fmt.Print("<> ")
			} else {
				fmt.Printf("%s ", s)
			}
		}
		fmt.Println()
	}

	for _, ln := range lines {
		put(ln)

		printTreeImpl(ln.Children, depth+1)
	}

}
