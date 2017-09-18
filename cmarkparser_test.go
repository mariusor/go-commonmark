package cmarkparser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type testPair struct {
	text     string
	expected bool
	doc      Document
}

type tests map[string]testPair

var emptyDoc = Document{}

func (n *NodeType) unmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("type should be a string, got %s", data)
	}

	got, ok := nodeTypeMap[s]
	if !ok {
		return fmt.Errorf("invalid NodeType %q", s)
	}
	*n = got
	return nil
}

func newNode(t NodeType, s string) Node {
	return Node{Type: t, Content: []byte(s)}
}

func newDoc(n Nodes) Document {
	return Document{
		Children: n,
	}
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
		newDoc(Nodes{newNode(Par, "some text")}),
	},
	// null char
	"null_char": {
		"\x00",
		true,
		newDoc(Nodes{newNode(Par, "\ufffd")}),
	},
	// spaces
	"space#1": {
		"\uc2a0",
		true,
		newDoc(Nodes{newNode(Par, "\uc2a0")}),
	},
	"space#2": {
		"\u2000",
		true,
		newDoc(Nodes{newNode(Par, "\u2000")}),
	},
	"space#3": {
		"\u2001",
		true,
		newDoc(Nodes{newNode(Par, "\u2001")}),
	},
	/*/
	// links, for now treated as paragraphs
	"link#1": {
		"[ana](httpslittrme)",
		true,
		newDoc(Nodes{newNode(Par, "[ana](httpslittrme)")}),
	},
	"link#2": {
		"[ana](https://littr.me)\n",
		true,
		newDoc(Nodes{newNode(Par, "[ana](https://littr.me)")}),
	},
	"link_after_text": {
		"some text before [test 123](https://littr.me)\n",
		true,
		newDoc(Nodes{newNode(Par, "some text before [test 123](https://littr.me)")}),
	},
	"link_before_text": {
		"[test 123](https://littr.me) some text after\n",
		true,
		newDoc(Nodes{newNode(Par, "[test 123](https://littr.me) some text after")}),
	},
	"link_inside_text": {
		"some text before [test 123](https://littr.me) some text after\n",
		true,
		newDoc(Nodes{newNode(Par, "some text before [test 123](https://littr.me) some text after")}),
	},
	/**/
	// utf8 only characters
	"utf8#1": {
		"𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș",
		true,
		newDoc(Nodes{newNode(Par, "𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș")}),
	},
	// thematic breaks
	"break#1:-": {
		" ---\n",
		true,
		newDoc(Nodes{newNode(TBreak, "-")}),
	},
	"break#2:*": {
		"  ***\n",
		true,
		newDoc(Nodes{newNode(TBreak, "*")}),
	},
	"break#3:*": {
		"  * * * *\n",
		true,
		newDoc(Nodes{newNode(TBreak, "*")}),
	},
	"break#4:-": {
		"   ___\r",
		true,
		newDoc(Nodes{newNode(TBreak, "_")}),
	},
	// misleading thematic break
	"not_a_break": {
		"   _*-*__",
		true,
		newDoc(Nodes{newNode(Par, "   _*-*__")}),
	},
	// headings
	"h1": {
		" # ana are mere\n",
		true,
		newDoc(Nodes{newNode(H1, "ana are mere")}),
	},
	"h2": {
		"## ana are mere\n",
		true,
		newDoc(Nodes{newNode(H2, "ana are mere")}),
	},
	"h3": {
		"  ### ana are mere\n",
		true,
		newDoc(Nodes{newNode(H3, "ana are mere")}),
	},
	"h4": {
		"#### ana are mere\n",
		true,
		newDoc(Nodes{newNode(H4, "ana are mere")}),
	},
	"h5": {
		"   #####  ana-are-mere\n",
		true,
		newDoc(Nodes{newNode(H5, "ana-are-mere")}),
	},
	"h6": {
		" ###### ana-are-mere\n",
		true,
		newDoc(Nodes{newNode(H6, "ana-are-mere")}),
	},
}

var trimb = func(s []byte) string {
	return strings.Trim(string(s), "\n\r")
}
var trims = func(s string) string {
	return strings.Trim(s, "\n\r")
}

func TestParse(t *testing.T) {

	var err error
	var doc Document
	for k, curTest := range someTests {
		t.Logf("Testing: %s", k)
		doc, err = Parse([]byte(curTest.text))

		if err != nil && curTest.expected {
			t.Errorf("Parse failed and success was expected %s\n %s", err, curTest.text)
		}
		if !reflect.DeepEqual(curTest.doc, doc) {
			t.Errorf("\n___\n%s\n___\n%s\n___", curTest.doc, doc)
		}
	}
}

func load_files(ext string) ([]string, error) {
	var files []string

	dir := "./tests"
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && path[len(path)-len(ext):] == ext {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not load files %s/*%s", dir, ext))
	}

	return files, nil
}

func get_file_contents(path string) []byte {
	f, _ := os.Open(path)

	data := make([]byte, 512)
	io.ReadFull(f, data)
	data = bytes.Trim(data, "\x00")

	return data
}

func testWithFiles(t *testing.T) {
	var tests []string
	var err error

	tests, err = load_files(".md")

	for _, path := range tests {
		var doc Document
		var res_doc Document
		data := get_file_contents(path)
		t.Logf("Testing: %s", path)
		res_path := fmt.Sprintf("%s.json", path[:len(path)-3])
		json.Unmarshal(get_file_contents(res_path), &res_doc)

		doc, err = Parse(data)

		if err == nil {
			log.Printf("%q", doc)
		}

		if err != nil {
			t.Errorf("%s", err)
		}
		if !reflect.DeepEqual(res_doc, doc) {
			t.Errorf("\n___\n%s\n___\n%s\n___", doc, res_doc)
		}

		t.Logf("%s", res_doc)
	}
}

func TestMain(m *testing.M) {
	if func(slice []string, s string) bool {
		for _, el := range slice {
			if s == el {
				return true
			}
		}
		return false
	}(os.Args, "quiet") {
		log.SetOutput(ioutil.Discard)
	}
	os.Exit(m.Run())
}
