package parser

import m "markdown"

func Parse(data []byte) (m.Document, error) {
	return parse(data)
}
