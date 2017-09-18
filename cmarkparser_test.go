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
	doc      testDocument
}

var emptyDoc = newDoc([]testNode{testNode{}}) //Document = Document{}

type testDocument struct {
	Children []testNode
}

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

type testNode struct {
	Type     NodeType
	Content  string
	Children []testNode `json:",omitempty"`
}

func newNode(t NodeType, s string) testNode {
	return testNode{Type: t, Content: s}
}

func newDoc(n []testNode) testDocument {
	return testDocument{
		Children: n,
	}
}

func (d *testDocument) String() string {
	var buffer bytes.Buffer
	for _, c := range d.Children {
		buffer.WriteString(fmt.Sprintf("  %s\n", c))
	}
	return buffer.String()
}

func (n *testNode) String() string {
	var buffer bytes.Buffer
	if len(n.Content) > 0 {
		buffer.WriteString(fmt.Sprintf("[%s] %s", n.Type, string(n.Content)))
	} else {
		buffer.WriteString(fmt.Sprintf("[%s]", n.Type))
	}
	if len(n.Children) > 0 {
		buffer.WriteString("\n[")
	}
	for _, c := range n.Children {
		buffer.WriteString(fmt.Sprintf("  %s\n", c))
	}
	if len(n.Children) > 0 {
		buffer.WriteString("]")
	}
	return buffer.String()
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
		newDoc([]testNode{newNode(Par, "some text")}),
	},
	// null char
	{
		"\x00",
		true,
		newDoc([]testNode{newNode(Par, "\ufffd")}),
	},
	// spaces
	{
		"\uc2a0",
		true,
		newDoc([]testNode{newNode(Par, "\uc2a0")}),
	},
	{
		"\u2000",
		true,
		newDoc([]testNode{newNode(Par, "\u2000")}),
	},
	{
		"\u2001",
		true,
		newDoc([]testNode{newNode(Par, "\u2001")}),
	},
	// links, for now treated as paragraphs
	{
		"[ana](httpslittrme)",
		true,
		newDoc([]testNode{newNode(Par, "[ana](httpslittrme)")}),
	},
	{
		"[ana](https://littr.me)\n",
		true,
		newDoc([]testNode{newNode(Par, "[ana](https://littr.me)")}),
	},
	{
		"some text before [test 123](https://littr.me)\n",
		true,
		newDoc([]testNode{newNode(Par, "some text before [test 123](https://littr.me)")}),
	},
	{
		"[test 123](https://littr.me) some text after\n",
		true,
		newDoc([]testNode{newNode(Par, "[test 123](https://littr.me) some text after")}),
	},
	{
		"some text before [test 123](https://littr.me) some text after\n",
		true,
		newDoc([]testNode{newNode(Par, "some text before [test 123](https://littr.me) some text after")}),
	},
	// utf8 only characters
	{
		"𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș",
		true,
		newDoc([]testNode{newNode(Par, "𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș")}),
	},
	// thematic breaks
	{
		" ---\n",
		true,
		newDoc([]testNode{newNode(TBreak, "-")}),
	},
	{
		"  ***\n",
		true,
		newDoc([]testNode{newNode(TBreak, "*")}),
	},
	{
		"  * * * *\n",
		true,
		newDoc([]testNode{newNode(TBreak, "*")}),
	},
	{
		"   ___\r",
		true,
		newDoc([]testNode{newNode(TBreak, "_")}),
	},
	// misleading thematic break
	{
		"   _*-*__",
		true,
		newDoc([]testNode{newNode(Par, "   _*-*__")}),
	},
	// headings
	{
		" # ana are mere\n",
		true,
		newDoc([]testNode{newNode(H1, "ana are mere")}),
	},
	{
		"## ana are mere\n",
		true,
		newDoc([]testNode{newNode(H2, "ana are mere")}),
	},

	{
		"  ### ana are mere\n",
		true,
		newDoc([]testNode{newNode(H3, "ana are mere")}),
	},
	{
		"#### ana are mere\n",
		true,
		newDoc([]testNode{newNode(H4, "ana are mere")}),
	},
	{
		"   #####  ana-are-mere\n",
		true,
		newDoc([]testNode{newNode(H5, "ana-are-mere")}),
	},
	{
		" ###### ana-are-mere\n",
		true,
		newDoc([]testNode{newNode(H6, "ana-are-mere")}),
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
	for _, curTest := range someTests {
		doc, err = Parse([]byte(curTest.text))

		if err != nil && curTest.expected {
			t.Errorf("Parse failed and success was expected %s\n %s", err, curTest.text)
		}
		if reflect.DeepEqual(curTest.doc, doc) {
			t.Errorf("Expected\n%s\ngot\n%s", curTest.doc, doc)
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

func TestWithFiles(t *testing.T) {
	var tests []string
	var err error

	tests, err = load_files(".md")

	for _, path := range tests {
		var doc Document
		var res_doc testDocument
		data := get_file_contents(path)
		log.Printf("%s:%s", path, path[:len(path)-3])
		res_path := fmt.Sprintf("%s.json", path[:len(path)-3])
		json.Unmarshal(get_file_contents(res_path), &res_doc)

		doc, err = Parse(data)

		if err == nil {
			log.Printf("%q", doc)
		}

		if err != nil {
			t.Errorf("%s", err)
		}
		if reflect.DeepEqual(res_doc, doc) {
			t.Errorf("\n____ expected ____\n%s\n______ got  ______\n%s", doc, res_doc)
		}
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
