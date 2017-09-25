package markdown

import (
	"bytes"
	"fmt"
	"strings"
)

type (
	NodeType   uint8
	Document   Node
	Nodes      []Node
	Attributes map[string]string
)

type Node struct {
	Type     NodeType
	Content  []byte
	Children Nodes
	//Attributes Attributes
}

const (
	None NodeType = iota
	Doc
	InlineText
	H1
	H2
	H3
	H4
	H5
	H6
	Par
	TBreak
)

var nodeTypeMap = map[string]NodeType{
	"nil": None,
	"doc": Doc,
	"txt": InlineText,
	"h1":  H1,
	"h2":  H2,
	"h3":  H3,
	"h4":  H4,
	"h5":  H5,
	"h6":  H6,
	"par": Par,
	"tbr": TBreak,
}

func NewNode() Node {
	return Node(Node{Type: None})
}

func NewDocument() Document {
	return Document(Node{Type: Doc})
}

func NewInlineText(cont []byte) Node {
	var el Node

	el.Type = InlineText
	el.Content = cont

	return el
}

func NewParagraph(cont []byte) Node {
	var el Node

	el.Type = Par
	el.Content = cont

	return el
}

func NewThematicBreak(t byte) Node {
	var el Node

	el.Type = TBreak
	el.Content = []byte{t}

	return el
}

func NewHeading(level uint, content []byte) Node {
	var el Node
	var t NodeType

	switch level {
	case 1:
		t = H1
	case 2:
		t = H2
	case 3:
		t = H3
	case 4:
		t = H4
	case 5:
		t = H5
	case 6:
		t = H6
	}
	el.Type = t
	el.Content = content

	return el
}

func (n *Node) Empty() bool {
	return n.Type == None /*|| len(n.Content) == 0*/
}

func (ns *Nodes) Empty() bool {
	return len(*ns) == 0
}

func (d *Document) Empty() bool {
	return len(d.Children) == 0
}

func (d Document) String() string {
	var r string = fmt.Sprintf("Document{%s}\n", d.Children)
	return r
}

func (n Node) String() string {
	var r string
	if len(n.Content) > 0 {
		r += fmt.Sprintf("[%s] \"%s\"", n.Type, n.Content)
	} else {
		r += fmt.Sprintf("[%s]", n.Type)
	}
	if len(n.Children) > 0 {
		r += fmt.Sprintf("\n\tChildren{%s}\n", n.Children)
	}
	return r
}

func (n NodeType) String() string {
	for key, node := range nodeTypeMap {
		if node == n {
			return key
		}
	}
	return "_nil"
}

func (n Nodes) String() string {
	var s string
	s += "{\n"
	for k, v := range n {
		s += fmt.Sprintf("\tNode{%d}: %s\n", k, v)
	}
	s += "}"
	return s
}

func GetNodeType(s string) NodeType {
	return nodeTypeMap[s]
}

func (d *Document) AddNodes(v interface{}) (bool, error) {
	if val, err := validNodeType(v); !val {
		return false, err
	}
	switch v.(type) {
	case Node:
		d.Children = append(d.Children, v.(Node))
	case Nodes:
		d.Children = append(d.Children, v.(Nodes)...)
	}
	return true, nil
}

func validNodeType(v interface{}) (bool, error) {
	switch t := v.(type) {
	case Node:
		//log.Printf("%s", t)
	case Nodes:
		//log.Printf("%s", t)
		//case Document, NodeType, Attributes:
	default:
		err := fmt.Errorf("invalid type '%s'", t)
		return false, err
	}
	return true, nil
}

func (n *Node) AddNodes(v interface{}) (bool, error) {
	if val, err := validNodeType(v); !val {
		return false, err
	}
	switch v.(type) {
	case Node:
		n.Children = append(n.Children, v.(Node))
	case Nodes:
		n.Children = append(n.Children, v.(Nodes)...)
	}

	return true, nil
}

func (n *Node) AppendContent(c []byte) (bool, error) {
	n.Content = append(n.Content, c...)
	return true, nil
}

var trimb = func(s []byte) []byte {
	return bytes.Trim(s, "\n\r ")
}
var trims = func(s string) string {
	return strings.Trim(s, "\n\r ")
}
