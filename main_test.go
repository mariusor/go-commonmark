package main

import (
	"./parser"
	"bytes"
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

var emptyDoc parser.Document = parser.Document{}

var someTests = []testPair{
	//    {
	//        "\u0000",
	//        true,
	//        parser.Document{
	//            Children: []parser.Node{
	//                parser.Node{
	//                    Type: parser.Par,
	//                    Content: []byte("\ufffd"),
	//                },
	//            },
	//        },
	//    },
	{
		"[ana](httpslittrme)",
		true,
		emptyDoc,
	},
	{
		"[ana](https://littr.me)",
		true,
		emptyDoc,
	},
	{
		"",
		false,
		emptyDoc,
	},
	{
		"some text before [test 123](https://littr.me)",
		true,
		emptyDoc,
	},
	{
		"[test 123](https://littr.me) some text after",
		true,
		emptyDoc,
	},
	{
		"some text before [test 123](https://littr.me) some text after",
		true,
		emptyDoc,
	},
	{
		"𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.Par,
					Content: []byte("𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș"),
				},
			},
		},
	},
	{
		" ---\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.TBreak,
					Content: []byte("-"),
				},
			},
		},
	},
	{
		"  ***\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.TBreak,
					Content: []byte("-"),
				},
			},
		},
	},
	{
		"  * * * *\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.TBreak,
					Content: []byte("-"),
				},
			},
		},
	},
	{
		"   ___\r",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.TBreak,
					Content: []byte("-"),
				},
			},
		},
	},
	{
		"   _*-*__\r",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.Par,
					Content: []byte("   _*-*__"),
				},
			},
		},
	},
	{
		"# ana are mere\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.H1,
					Content: []byte("ana are mere"),
				},
			},
		},
	},
	{
		"## ana are mere\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.H2,
					Content: []byte("ana are mere"),
				},
			},
		},
	},

	{
		"### ana are mere\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.H3,
					Content: []byte("ana are mere"),
				},
			},
		},
	},
	{
		"#### ana are mere\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.H4,
					Content: []byte("ana are mere"),
				},
			},
		},
	},
	{
		"##### ana-are-mere\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.H5,
					Content: []byte("ana-are-mere"),
				},
			},
		},
	},
	{
		"###### ana-are-mere\n",
		true,
		parser.Document{
			Children: []parser.Node{
				parser.Node{
					Type:    parser.H6,
					Content: []byte("ana-are-mere"),
				},
			},
		},
	},
}

func TestParse(t *testing.T) {
	//    var n parser.Node
	//    t.Logf("%s", n.String())

	for _, curTest := range someTests {
		doc, err := parser.Parse([]byte(curTest.text))
		t.Logf("Testing %q", strings.Trim(curTest.text, "\n\r"))

		if err != nil && curTest.expected {
			t.Errorf("Parse result invalid, expected '%t, got %v'\n", curTest.expected)
		}

		if !curTest.doc.Equal(doc) {
			t.Errorf("Expected %q\nGot %q", curTest.doc.String(), doc.String())
		}

		testChildren := curTest.doc.Children
		children := doc.Children
		if len(testChildren) != len(children) {
			t.Errorf("Parse result invalid, children length expected %d != %d", len(testChildren), len(children))
		}

		if len(testChildren) > 0 && len(children) > 0 {
			testNode := testChildren[0]
			node := children[0]
			if testNode.Type != node.Type {
				t.Errorf("Parse result invalid, node type expected %v != %v", testNode.Type, node.Type)
			}
			if !bytes.Equal(testNode.Content, node.Content) {
				t.Errorf("Parse result invalid, node content expected %s != %s", testNode.Content, node.Content)
				t.Errorf("Parse result invalid, node content expected %v != %v", testNode.Content, node.Content)
			}
			t.Logf("%s", doc.String())
		}
	}
}

func TestParseReadme(t *testing.T) {
	f, _ := os.Open("README.md")

	data := make([]byte, 512)
	io.ReadFull(f, data)
	data = bytes.Trim(data, "\x00")

	if _, err := parser.Parse(data); err != nil {
		t.Errorf("\tParse invalid\n")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
