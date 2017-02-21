package cmarkparser

import (
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
	doc      Document
}

var emptyDoc = newDoc([]Node{Node{}}) //Document = Document{}

func newDoc(n []Node) Document {
	return Document{
		Children: n,
	}
}

func newNode(t NodeType, s string) Node {
	return Node{Type: t, Content: []byte(s)}
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
		newDoc([]Node{newNode(Par, "some text")}),
	},
	// null char
	{
		"\u0000\n",
		true,
		newDoc([]Node{newNode(Par, "\xff\xfd")}),
	},
	// spaces
	{
		"\uc2a0",
		true,
		newDoc([]Node{newNode(Par, "\uc2a0")}),
	},
	{
		"\u2000",
		true,
		newDoc([]Node{newNode(Par, "\u2000")}),
	},
	{
		"\u2001",
		true,
		newDoc([]Node{newNode(Par, "\u2001")}),
	},
	// links, for now treated as paragraphs
	//	{
	//		"[ana](httpslittrme)",
	//		true,
	//		newDoc([]Node{newNode(Par, "[ana](httpslittrme)")}),
	//	},
	//	{
	//		"[ana](https://littr.me)\n",
	//		true,
	//		newDoc([]Node{newNode(Par, "[ana](https://littr.me)")}),
	//	},
	//	{
	//		"some text before [test 123](https://littr.me)\n",
	//		true,
	//		newDoc([]Node{newNode(Par, "some text before [test 123](https://littr.me)")}),
	//	},
	//	{
	//		"[test 123](https://littr.me) some text after\n",
	//		true,
	//		newDoc([]Node{newNode(Par, "[test 123](https://littr.me) some text after")}),
	//	},
	//	{
	//		"some text before [test 123](https://littr.me) some text after\n",
	//		true,
	//		newDoc([]Node{newNode(Par, "some text before [test 123](https://littr.me) some text after")}),
	//	},
	// utf8 only characters
	{
		"𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș",
		true,
		newDoc([]Node{newNode(Par, "𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș")}),
	},
	// thematic breaks
	{
		" ---\n",
		true,
		newDoc([]Node{newNode(TBreak, "-")}),
	},
	{
		"  ***\n",
		true,
		newDoc([]Node{newNode(TBreak, "*")}),
	},
	{
		"  * * * *\n",
		true,
		newDoc([]Node{newNode(TBreak, "*")}),
	},
	{
		"   ___\r",
		true,
		newDoc([]Node{newNode(TBreak, "_")}),
	},
	// misleading thematic break
	{
		"   _*-*__",
		true,
		newDoc([]Node{newNode(Par, "   _*-*__")}),
	},
	// headings
	{
		" # ana are mere\n",
		true,
		newDoc([]Node{newNode(H1, "ana are mere")}),
	},
	{
		"## ana are mere\n",
		true,
		newDoc([]Node{newNode(H2, "ana are mere")}),
	},

	{
		"  ### ana are mere\n",
		true,
		newDoc([]Node{newNode(H3, "ana are mere")}),
	},
	{
		"#### ana are mere\n",
		true,
		newDoc([]Node{newNode(H4, "ana are mere")}),
	},
	{
		"   #####  ana-are-mere\n",
		true,
		newDoc([]Node{newNode(H5, "ana-are-mere")}),
	},
	{
		" ###### ana-are-mere\n",
		true,
		newDoc([]Node{newNode(H6, "ana-are-mere")}),
	},
}

var readmeTest = func() testPair {
	f, _ := os.Open("README.md")

	data := make([]byte, 512)
	io.ReadFull(f, data)
	data = bytes.Trim(data, "\x00")

	title := newNode(H1, "Ragel playground")
	hr := newNode(TBreak, "-")
	p1 := newNode(Par, "A small go repository to learn some ragel usage by implementing a Common Mark ")
	p2 := newNode(Par, "Using the [0.27](http://spec.commonmark.org/0.27/) version of the specification.")
	p3 := newNode(Par, "[![Build Status](https://travis-ci.org/mariusor/ragel-playgrnd.svg?branch=master)](https://travis-ci.org/mariusor/ragel-playgrnd)")
	d := newDoc([]Node{title, hr, p1, p2, p3})

	return testPair{
		text:     string(data),
		expected: true,
		doc:      d,
	}
}

var trimb = func(s []byte) string {
	return strings.Trim(string(s), "\n\r")
}
var trims = func(s string) string {
	return strings.Trim(s, "\n\r")
}

func assertDocumentsEqual(d1 Document, d2 Document) (bool, error) {
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

func assertNodesEqual(n1 Node, n2 Node) (bool, error) {
	if n1.Type != n2.Type {
		return false, errors.New(fmt.Sprintf("  Node type expected %q != %q", n1.Type.String(), n2.Type.String()))
	}
	if !bytes.Equal(n1.Content, n2.Content) {
		return false, errors.New(fmt.Sprintf("  Node content expected %q:%v != %q:%v", trimb(n1.Content), n1.Content, trimb(n2.Content), n2.Content))
	}
	return true, nil
}

func TestParse(t *testing.T) {
	//someTests = append(someTests, readmeTest())

	var err error
	var doc Document
	for _, curTest := range someTests {
		doc, err = Parse([]byte(curTest.text))

		//if err != nil && curTest.expected {
		//	t.Errorf(" Parse result invalid, expected %t, got %v\n", curTest.expected)
		//}
		_, err = assertDocumentsEqual(curTest.doc, doc)
		if err != nil {
			t.Errorf("\n%s", err)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
