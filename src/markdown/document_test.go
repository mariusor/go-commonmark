package markdown

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var stopOnFailure = false

func TestDocument_Empty(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	d := NewDocument()

	if !d.Empty() {
		err_f("New document should be empty")
	}
}

func TestNewDocument(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	d := NewDocument()

	if d.Type != Doc {
		err_f("Invalid document type %s", d.Type)
	}

	if len(d.Children) > 0 {
		err_f("Invalid document children count %d > %d", len(d.Children), 0)
	}

	if len(d.Content) > 0 {
		err_f("Invalid document content. Should be empty: %s", string(d.Content))
	}
}

func TestNewEmptyNode(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	e := NewEmptyNode()

	if e.Type != None {
		err_f("Invalid node type %s", e.Type)
	}

	if len(e.Children) > 0 {
		err_f("Invalid node children count %e > %e", len(e.Children), 0)
	}

	if len(e.Content) > 0 {
		err_f("Invalid node content. Should be empty: %s", string(e.Content))
	}
}

func TestNewHeading(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	h1 := NewHeading(1, []byte{})
	if h1.Type != H1 {
		err_f("Invalid heading type %s", h1.Type)
	}
	h2 := NewHeading(2, []byte{})
	if h2.Type != H2 {
		err_f("Invalid heading type %s", h2.Type)
	}
	h3 := NewHeading(3, []byte{})
	if h3.Type != H3 {
		err_f("Invalid heading type %s", h3.Type)
	}
	h4 := NewHeading(4, []byte{})
	if h4.Type != H4 {
		err_f("Invalid heading type %s", h4.Type)
	}
	h5 := NewHeading(5, []byte{})
	if h5.Type != H5 {
		err_f("Invalid heading type %s", h5.Type)
	}
	h6 := NewHeading(6, []byte{})
	if h6.Type != H6 {
		err_f("Invalid heading type %s", h6.Type)
	}
}

func TestNewInlineText(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	cont := []byte{}
	tx := NewInlineText(cont)
	if tx.Type != InlineText {
		err_f("Invalid text type %s", tx.Type)
	}
	if string(tx.Content) != string(cont) {
		err_f("Invalid content %s != %s", tx.Type, cont)
	}
}

func TestNewParagraph(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	cont := []byte{}
	p := NewParagraph(cont)
	if p.Type != Par {
		err_f("Invalid text type %s", p.Type)
	}
	if string(p.Content) != string(cont) {
		err_f("Invalid content %s != %s", p.Type, cont)
	}
}

func TestNewThematicBreak(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	var cont byte = 0x23
	br := NewThematicBreak(cont)
	if br.Type != TBreak {
		err_f("Invalid text type %s", br.Type)
	}
}

func TestNode_Empty(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	d := NewEmptyNode()

	if d.Type != None {
		err_f("Invalid node type %s", d.Type)
	}

	if !d.Empty() {
		err_f("Node should be empty")
	}

	if len(d.Children) > 0 {
		err_f("Invalid node children count %d > %d", len(d.Children), 0)
	}

	if len(d.Content) > 0 {
		err_f("Invalid node content. Should be empty: %s", string(d.Content))
	}
}

func TestDocument_String(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}

	e := NewDocument()
	e_s := "Document:{{\n}}\n"
	if e_s != e.String() {
		err_f("Empty document string invalid: \n%s\n%s", e.String(), e_s)
	}
}

func TestNode_String(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}
	e := NewEmptyNode()
	e_s := "[nil]"

	if e.String() != e_s {
		err_f("Empty node string invalid: \n%s\n%s", e.String(), e_s)
	}
}

func TestNodes_String(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}

	e := NewEmptyNode()
	n := Nodes{e}
	n_s := fmt.Sprintf("{\n\tNode{0}: %s\n}", e)
	if n_s != n.String() {
		err_f("Empty nodes string invalid: \n%s\n%s", n.String(), n_s)
	}
}

func TestNodeType_String(t *testing.T) {
	err_f := t.Errorf
	if stopOnFailure {
		err_f = t.Fatalf
	}

	for nt_s, nt := range nodeTypeMap {
		if nt_s != nt.String() {
			err_f("Node type string invalid: \n%s\n%s", nt.String(), nt_s)
		}
	}

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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	os.Exit(m.Run())
}
