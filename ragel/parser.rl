// -*-go-*-
//
// Commonmark parser
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
//

package parser

import(
    "bytes"
    "errors"
    "log"

    m "github.com/mariusor/go-commonmark/src/markdown"
)

func arr_splice(dst []byte, src []byte, pos int) []byte {
    var ret = make([]byte, 0)
    for _, a := range dst[:pos] {
        ret = append(ret, a)
    }
    for _, b := range src {
        ret = append(ret, b)
    }
    for _, c := range dst[pos+1:] {
        ret = append(ret, c)
    }
    return ret
}

%% machine parser;
%% include character_definitions "characters.rl";
%% include blocks "blocks.rl";

%% write data;

func parse(data []byte) (m.Document, error) {
    //t_data = trimb(data)
    cs, p, pe := 0, 0, len(data)
    ts, te, act := 0, 0, 0
    stack := make([]int,20)
    top := 0
    eof := len(data)

    var node m.Node
    var nodes m.Nodes;
    var doc m.Document = m.NewDocument()

    var heading_level uint;
    var thematic_break_symbol byte
    var end_of_heading int = 0
    var end_of_par int = 0
    if false {
        log.Printf("ts:%d", ts)
        log.Printf("sym: %s lvl: %d", string(thematic_break_symbol), heading_level)
    }

    //fmt.Printf("%s", data)
    if pe == 0 {
        return doc, errors.New("Empty document")
    }

    var mark int
    %%{
        action emit_eof {
            if (len(nodes) > 0) {
                doc.AddNodes(nodes)
            }
            log.Printf("emit_end_of_file:(%d)", eof)
        }

        action emit_add_block {
            if !node.Empty() {
                nodes = append(nodes, node)
                log.Printf("emit_add_block(%d): %#v", p, node)
                node = m.NewNode()
            }
        }

        action emit_add_thematic_break {
            nodes = append(nodes, node)
            log.Printf("emit_add_thematic_break(%d) %#v", p, node)
            node = m.NewNode()
        }

        action emit_add_atx_heading {
            nodes = append(nodes, node)
            log.Printf("emit_add_atx_heading(%d) %#v", p, node)
            node = m.NewNode()
        }

#        document := |*
#           text_paragraph => emit_add_block;
#           thematic_break => emit_add_thematic_break;
#           atx_heading => emit_add_atx_heading;
#       *|;

        main := block %emit_add_block %eof emit_eof;

        write init;
        write exec;
    }%%

    if false {
        //log.Printf("end_of_par %d", end_of_par)
        log.Printf("stack %d %v", top, stack)
        //log.Printf("eoh %d", end_of_heading)
        log.Printf("last node %s", node)
        log.Printf("nodes %s", nodes)
        log.Printf("doc %s", doc)
        log.Printf("mark:%d, te:%d, act:%d", mark, te, act)
    }
    return doc, nil
}
