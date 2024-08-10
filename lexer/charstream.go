package lexer

const eof_char rune = 0

type CharStream struct {
	text    []byte
	idx     int
	current rune

	lineIdx int
}

func NewCharStream(txt []byte) CharStream {
	return CharStream{txt, 0, 0, 0}
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
	r := c.getAt(c.idx)
	c.idx++
	return r
}

func (c *CharStream) getAt(i int) rune {
	if i < 0 || i >= len(c.text) {
		return eof_char
	}
	return rune(c.text[i])
}
