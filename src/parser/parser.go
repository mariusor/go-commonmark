package parser

import m "github.com/mariusor/go-commonmark/src/markdown"

func Parse(data []byte) (m.Document, error) {
	return parse(data)
}
