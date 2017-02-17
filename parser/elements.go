package parser

import (
    "fmt"
    "bytes"
)

func NewParagraph(cont []byte) Node {
    var el Node

    el.Type = Par
    el.Content = cont

    return el
}

func NewThematicBreak(t byte) Node {
    var el Node

    el.Type = TBreak
    el.Content = append(el.Content, t)

    return el
}

func NewHeader(level uint, content []byte) Node {
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
        t= H6
    }
    el.Type = t
    el.Content = content

    return el
}

func (n *Node)Equal(tn Node) bool {
    return true
}

func (d *Document)Equal(td Document) bool {
    return true
}

func (d *Document) String() string {
    var buffer bytes.Buffer
    for _, c := range(d.Children) {
        buffer.WriteString(c.String())
    }
    return buffer.String()
}

func (n *Node) String() string {
    var buffer bytes.Buffer
    buffer.WriteString(fmt.Sprintf("[%s] %s\n", n.Type.String(), string(n.Content)))
    for _, c := range(n.Children) {
        buffer.WriteString(c.String())
    }
    return buffer.String()
}


type NodeType uint8

const (
    None NodeType = iota 
    Doc
    H1
    H2
    H3
    H4
    H5
    H6
    Par
    TBreak
)

func (nt *NodeType) String() string {
    switch *nt {
    case Doc:
        return "document"
    case Par:
        return "paragraph"
    case TBreak:
        return "thematic break"
    case H1:
       return "h1"
    case H2:
       return "h2"
    case H3:
       return "h3"
    case H4:
       return "h4"
    case H5:
       return "h5"
    case H6:
       return "h6"
   default:
       return "none"
    }
}

type Document struct {
    Children []Node
}

type Node struct {
   Type NodeType
   Content []byte
   Children []Node
   //Attributes map[string]string
}
