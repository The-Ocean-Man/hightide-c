package main

import (
	"fmt"
	"log"
	"os"

	"github.com/The-Ocean-Man/hightide-c/lexer"
)

func main() {
	bytes, err := os.ReadFile("./test.ht")
	// fmt.Println(bytes)

	if err != nil {
		log.Fatal(err)
	}

	stream := lexer.NewCharStream(bytes)
	// b := make([]rune, 0)

	// for {
	// 	r := stream.Next()
	// 	b = append(b, r)

	// 	if r == 0 {
	// 		break
	// 	}
	// }

	// fmt.Println(b)

	// return

	l := lexer.NewLexer(&stream)

	lines := l.Parse()
	fmt.Println(lines[0].Children[0])
	// for _, ln := range lines {
	// 	printTree(ln)
	// }
}
