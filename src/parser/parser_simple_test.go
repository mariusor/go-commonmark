package parser

import (
	"io/ioutil"
	"log"
	m "markdown"
	"os"
	"reflect"
	"testing"
)

var stopOnFailure = false

type testPair struct {
	text     string
	expected bool
	doc      m.Document
}

type tests map[string]testPair

var emptyDoc = m.NewDocument()

func newNode(t m.NodeType, s string, c m.Nodes) m.Node {
	node := m.NewNode()
	node.Type = t
	node.Content = []byte(s)
	node.Children = c
	return node
}

func newDoc(n m.Nodes) m.Document {
	doc := m.NewDocument()
	doc.Type = m.Doc
	doc.Children = n
	return doc
}

var someTests = tests{
	// empty doc
	"empty": {
		"",
		false,
		emptyDoc,
	},
	"line": {
		"some text",
		true,
		newDoc(m.Nodes{newNode(m.Par, "some text", nil)}),
	},
	"line_eol": {
		"some text\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "some text", nil)}),
	},
	"line_mac_eol": {
		"some text\r",
		true,
		newDoc(m.Nodes{newNode(m.Par, "some text", nil)}),
	},
	"line_win_eol": {
		"some text\r\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "some text", nil)}),
	},
	//"2line_eol": {
	//	"some text\ntest\n",
	//	true,
	//	newDoc(m.Nodes{newNode(m.Par, "some text\ntest", nil)}),
	//},
	"line_eop": {
		"some text\n\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "some text", nil)}),
	},
	// null char
	"null_char": {
		"\x00",
		true,
		newDoc(m.Nodes{newNode(m.Par, "\ufffd", nil)}),
	},
	// spaces
	"space#1": {
		"\uc2a0",
		true,
		newDoc(m.Nodes{newNode(m.Par, "\uc2a0", nil)}),
	},
	"space#2": {
		"\u2000",
		true,
		newDoc(m.Nodes{newNode(m.Par, "\u2000", nil)}),
	},
	"space#3": {
		"\u2001",
		true,
		newDoc(m.Nodes{newNode(m.Par, "\u2001", nil)}),
	},
	// links, for now treated as paragraphs
	"link#1": {
		"[ana](httpslittrme)",
		true,
		newDoc(m.Nodes{newNode(m.Par, "[ana](httpslittrme)", nil)}),
	},
	"link#2": {
		"[ana](https://littr.me)\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "[ana](https://littr.me)", nil)}),
	},
	"link_after_text": {
		"some text before [test 123](https://littr.me)\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "some text before [test 123](https://littr.me)", nil)}),
	},
	"link_before_text": {
		"[test 123](https://littr.me) some text after\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "[test 123](https://littr.me) some text after", nil)}),
	},
	"link_inside_text": {
		"some text before [test 123](https://littr.me) some text after\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "some text before [test 123](https://littr.me) some text after", nil)}),
	},
	// utf8 only characters
	"utf8#1": {
		"𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș",
		true,
		newDoc(m.Nodes{newNode(m.Par, "𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș", nil)}),
	},
	// thematic breaks
	"break#1:-": {
		" ---\n\n",
		true,
		newDoc(m.Nodes{newNode(m.TBreak, "-", nil)}),
	},
	"break#2:*": {
		"  ***\n\n",
		true,
		newDoc(m.Nodes{newNode(m.TBreak, "*", nil)}),
	},
	"break#3:*": {
		"  * * * *\n\n",
		true,
		newDoc(m.Nodes{newNode(m.TBreak, "*", nil)}),
	},
	"break#4:-": {
		"   ___\r\r",
		true,
		newDoc(m.Nodes{newNode(m.TBreak, "_", nil)}),
	},
	// misleading thematic break
	"not_a_break": {
		"   _*-*__\r\n",
		true,
		newDoc(m.Nodes{newNode(m.Par, "   _*-*__", nil)}),
	},
	// headings
	"h1": {
		" # ana are mere\n\n",
		true,
		newDoc(m.Nodes{newNode(m.H1, "ana are mere", nil)}),
	},
	"h2": {
		"## ana are mere\n",
		true,
		newDoc(m.Nodes{newNode(m.H2, "ana are mere", nil)}),
	},
	"h3": {
		"  ### ana are mere\n",
		true,
		newDoc(m.Nodes{newNode(m.H3, "ana are mere", nil)}),
	},
	"h4": {
		"#### ana are mere\n",
		true,
		newDoc(m.Nodes{newNode(m.H4, "ana are mere", nil)}),
	},
	"h5": {
		"   #####  ana-are-mere\n",
		true,
		newDoc(m.Nodes{newNode(m.H5, "ana-are-mere", nil)}),
	},
	"h6": {
		" ###### ana-are-mere\n",
		true,
		newDoc(m.Nodes{newNode(m.H6, "ana-are-mere", nil)}),
	},
	"h1_with_end": {
		" # ana are mere #\n",
		true,
		newDoc(m.Nodes{newNode(m.H1, "ana are mere", nil)}),
	},
	"h2_with_end": {
		"## ana are mere ##\n",
		true,
		newDoc(m.Nodes{newNode(m.H2, "ana are mere", nil)}),
	},
	"h3_with_mismatched_end": {
		"### ana are mere ##\n",
		true,
		newDoc(m.Nodes{newNode(m.H3, "ana are mere", nil)}),
	},
}

func TestSimpleParse(t *testing.T) {
	var err error
	var doc m.Document
	var errF = t.Errorf
	if stopOnFailure {
		errF = t.Fatalf
	}
	for k, curTest := range someTests {
		t.Logf("Testing: %s", k)
		doc, err = Parse([]byte(curTest.text))

		if err != nil && curTest.expected {
			errF("Parse failed and success was expected %s\n %s", err, curTest.text)
		}
		if !reflect.DeepEqual(doc, curTest.doc) {
			errF("\n%#v\n%#v\n%s\n%s", doc, curTest.doc, doc, curTest.doc)
		}
	}
}

type testDoc struct {
	Children []testNode
}

type testNode struct {
	Type     string
	Content  []string
	Children []testNode
}

func (t *testDoc) Document() m.Document {
	d := m.NewDocument()
	for _, v := range t.Children {
		d.Children = append(d.Children, *v.Node())
	}
	return d
}

func (t *testNode) Node() *m.Node {
	n := m.Node{}
	for _, v := range t.Children {
		n.Children = append(n.Children, *v.Node())
	}
	for k, s := range t.Content {
		if k > 0 {
			n.Content = append(n.Content, byte('\n'))
		}
		for _, b := range s {
			n.Content = append(n.Content, byte(b))
		}
	}
	n.Type = m.GetNodeType(t.Type)
	return &n
}

func TestMain(m *testing.M) {
	f := func(slice []string, s string) bool {
		for _, el := range slice {
			if s == el {
				return true
			}
		}
		return false
	}
	if f(os.Args, "stop-on-fail") {
		stopOnFailure = true
	}
	if f(os.Args, "quiet") {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile & 0)
	os.Exit(m.Run())
}
