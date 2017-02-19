// -*-go-*-
//
// Commonmark link parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
// MIT License
// 

package parser

import(
    "fmt"
    "errors"
)

func Parse (data []byte) (Document, error) {
    return parse(data)
}

%% machine parser;
%% include character_definitions "characters.rl";
%% include blocks "blocks.rl";

%% write data;

func parse(data []byte) (Document, error) {
    cs, p, pe := 0, 0, len(data)
    eof := len(data)

    var doc Document = Document{Children: []Node{Node{}}}
    if pe == 0 {
        return doc, errors.New("Empty document")
    }

    var node Node
    var header_level uint;
    var nodes []Node;
    //fmt.Printf("Incoming str: %#v - len %d\n", data, len(data))

    var mark int

    %%{
        action emit_eof {
            if !node.Empty() {
                nodes = append(nodes, node)
            }
            doc.Children = nodes
        }

        action emit_add_block {
            if !node.Empty() {
                nodes = append(nodes, node)
            }
            node = Node{}
        }

        document = (block %emit_add_block)*;
        main := document %eof emit_eof;
 
        write init;
        write exec;
    }%%

    //fmt.Printf("last node: %v\n", node)
    //fmt.Printf("last mark: %v\n", mark)
    //fmt.Printf("last level: %v\n", header_level)

    fmt.Printf("")

    return doc, nil
}
