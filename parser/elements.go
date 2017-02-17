package parser

import (
//    "fmt"
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

type NodeType uint8

const (
    None NodeType = iota 
    H1
    H2
    H3
    H4
    H5
    H6
    Par
    TBreak
)

type Document struct {
    Children []Node
}

type Node struct {
   Type NodeType
   Content []byte
   //Children []Node
   //Attributes map[string]string
}
