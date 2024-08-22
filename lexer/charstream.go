package lexer

import (
	"bufio"
)

const eof_char rune = '\x00'

type CharStream struct {
	// text    *string
	reader  *bufio.Reader
	idx     int
	current rune

	lineIdx  int
	firstEOF bool // Returns \n once when reached eof
}

func NewCharStream(r *bufio.Reader) CharStream {
	return CharStream{r, 0, 0, 0, true}
}

func (c *CharStream) Next() rune {
	r := c.get()
	// fmt.Printf("Nexty %c\n", r)

	if r == '\r' {
		// fmt.Println(r)
		return c.Next()
	}

	if r == '\n' {
		c.lineIdx++
	}

	c.current = r
	return r
}

func (c *CharStream) Current() rune {
	return c.current
}

func (c *CharStream) get() rune {
	r, _, err := c.reader.ReadRune()

	if err != nil {
		if c.current != '\n' && c.firstEOF {
			c.firstEOF = false
			return '\n'
		}

		return eof_char
	}

	return r
}
