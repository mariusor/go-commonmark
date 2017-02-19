package main

import (
	"./parser"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type testPair struct {
	text     string
	expected bool
	doc      parser.Document
}

var emptyDoc = newDoc([]parser.Node{parser.Node{}}) //parser.Document = parser.Document{}

func newDoc(n []parser.Node) parser.Document {
	return parser.Document{
		Children: n,
	}
}

func newNode(t parser.NodeType, s string) parser.Node {
	return parser.Node{Type: t, Content: []byte(s)}
}

var someTests = []testPair{
	// empty doc
	{
		"",
		false,
		emptyDoc,
	},
	{
		"some text",
		true,
		newDoc([]parser.Node{newNode(parser.Par, "some text")}),
	},
	// null char
	//	{
	//		"\u0000\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte{0x3f},
	//				},
	//			},
	//		},
	//	},
	// spaces
	//	{
	//		"\uc2a0\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("\uc2a0"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"\u2000\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("\u2000"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"\u2001\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("\u2001"),
	//				},
	//			},
	//		},
	//	},
	//	// links, for now treated as paragraphs
	//	{
	//		"[ana](httpslittrme)\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("[ana](httpslittrme)"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"[ana](https://littr.me)\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("[ana](https://littr.me)"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"some text before [test 123](https://littr.me)\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("some text before [test 123](https://littr.me)"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"[test 123](https://littr.me) some text after\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("[test 123](https://littr.me) some text after"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"some text before [test 123](https://littr.me) some text after\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("some text before [test 123](https://littr.me) some text after"),
	//				},
	//			},
	//		},
	//	},
	//	// utf8 only characters
	//	{
	//		"𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș"),
	//				},
	//			},
	//		},
	//	},
	//	// headings
	//	{
	//		" ---\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.TBreak,
	//					Content: []byte("-"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"  ***\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.TBreak,
	//					Content: []byte("-"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"  * * * *\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.TBreak,
	//					Content: []byte("-"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"   ___\r",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.TBreak,
	//					Content: []byte("-"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"   _*-*__\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.Par,
	//					Content: []byte("   _*-*__"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"# ana are mere\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.H1,
	//					Content: []byte("ana are mere"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"## ana are mere\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.H2,
	//					Content: []byte("ana are mere"),
	//				},
	//			},
	//		},
	//	},
	//
	//	{
	//		"### ana are mere\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.H3,
	//					Content: []byte("ana are mere"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"#### ana are mere\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.H4,
	//					Content: []byte("ana are mere"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		"#####  ana-are-mere\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.H5,
	//					Content: []byte("ana-are-mere"),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		" ###### ana-are-mere\n",
	//		true,
	//		parser.Document{
	//			Children: []parser.Node{
	//				parser.Node{
	//					Type:    parser.H6,
	//					Content: []byte("ana-are-mere"),
	//				},
	//			},
	//		},
	//	},
}
var trimb = func(s []byte) string {
	return strings.Trim(string(s), "\n\r")
}
var trims = func(s string) string {
	return strings.Trim(s, "\n\r")
}

func assertDocumentsEqual(d1 parser.Document, d2 parser.Document) (bool, error) {
	if !d1.Equal(d2) {
		return false, errors.New(fmt.Sprintf("Expected %q, got %q", trims(d1.String()), trims(d2.String())))
	}
	d1Children := d1.Children
	d2Children := d2.Children
	if len(d1Children) != len(d2Children) {
		return false, errors.New(fmt.Sprintf(" Children length expected %d != %d", len(d1Children), len(d2Children)))
	}
	if len(d1Children) > 0 && len(d2Children) > 0 {
		//t.Logf("%s", dt.String())
		for i, n1 := range d1Children {
			status, err := assertNodesEqual(n1, d2Children[i])
			if err != nil {
				return status, err
			}
		}
	}
	return true, nil
}

func assertNodesEqual(n1 parser.Node, n2 parser.Node) (bool, error) {
	if n1.Type != n2.Type {
		return false, errors.New(fmt.Sprintf("  Node type expected %q != %q", n1.Type.String(), n2.Type.String()))
	}
	if !bytes.Equal(n1.Content, n2.Content) {
		return false, errors.New(fmt.Sprintf("  Node content expected %q:%v != %q:%v", trimb(n1.Content), n1.Content, trimb(n2.Content), n2.Content))
	}
	return true, nil
}

func TestParse(t *testing.T) {
	for _, curTest := range someTests {
		doc, err := parser.Parse([]byte(curTest.text))

		//t.Logf("Testing %q", trims(curTest.text))

		if err != nil && curTest.expected {
			t.Errorf(" Parse result invalid, expected %t, got %v\n", curTest.expected)
		}

		_, errs := assertDocumentsEqual(curTest.doc, doc)
		if errs != nil {
			t.Errorf("%s", errs)
		}
	}
}

func TestParseReadme(t *testing.T) {
	f, _ := os.Open("README.md")

	data := make([]byte, 512)
	io.ReadFull(f, data)
	data = bytes.Trim(data, "\x00")

	t.Logf("\n%s\n", data)
	ast, err := parser.Parse(data)
	if err != nil {
		t.Errorf("\tParse invalid\n")
	}

	title := newNode(parser.H1, "Ragel playground")
	hr := newNode(parser.TBreak, "-")
	p1 := newNode(parser.Par, "A small go repository to learn some ragel usage by implementing a Common Mark parser.")
	p2 := newNode(parser.Par, "Using the [0.27](http://spec.commonmark.org/0.27/) version of the specification.")
	p3 := newNode(parser.Par, "[![Build Status](https://travis-ci.org/mariusor/ragel-playgrnd.svg?branch=master)](https://travis-ci.org/mariusor/ragel-playgrnd)")
	doc := newDoc([]parser.Node{title, hr, p1, p2, p3})

	t.Logf("%q\n%q\n", doc.String(), ast.String())
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
