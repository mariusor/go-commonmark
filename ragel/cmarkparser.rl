// -*-go-*-
//
// Commonmark parser
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
//

package cmarkparser

import(
    "log"
    "bytes"
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
    //t_data = trimb(data)
    cs, p, pe := 0, 0, len(data)
    ts, te, act := 0, 0, 0
    log.Printf("ts:%d", ts)
    eof := len(data)

    var node Node
    var doc Document
    var heading_level uint;
    var thematic_break_symbol byte
    var nodes Nodes;

    if pe == 0 {
        return doc, errors.New("Empty document")
    }
    doc = NewDocument()
    node = Node{}
    //fmt.Printf("Incoming[%d]: \"%s\"", len(data), data)

    var mark int
    //var content []byte

    %%{
        action emit_eof {
            if doc.Empty() {
                node = NewParagraph(data[:p])
                //log.Printf("current node: %s\n", node)
            }
            if len(nodes) == 0 {
                nodes = append(nodes, node)
            }
            if (len(nodes) > 0) {
                doc.Children = nodes
            }
            log.Printf("emit_end_of_file:(%d)", eof)
        }

        action emit_add_block {
            if !node.Empty() && node.Type != Doc {
                nodes = append(nodes, node)
                //log.Printf("appending node: %s\n", node)
                node = NewEmptyNode()
            }
            log.Printf("emit_add_block(%d)", p)
        }

        #single_line_doc = line_char* (eop | eol)? >emit_add_line;
        document = (block %emit_add_block)*;

#        main := |*
#            block => emit_add_block;
#            single_line_doc => emit_add_line;
#        *|;

        main := document %eof(emit_eof);

        write init;
        write exec;
    }%%

    //log.Printf("last node %s", node)
    //log.Printf("nodes %s", nodes)
    //log.Printf("doc %s", doc)
    log.Printf("mark:%d, te:%d, act:%d", mark, te, act)
    return doc, nil
}
