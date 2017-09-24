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
    "fmt"
    "log"
    m "markdown"
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
    eof := len(data)

    var node m.Node
    var doc m.Document
    var nodes m.Nodes;

    var heading_level uint;
    var thematic_break_symbol byte
    if false {
        log.Printf("ts:%d", ts)
        log.Printf("sym: %s lvl: %d", string(thematic_break_symbol), heading_level)
    }

    fmt.Printf("%s", data)
    if pe == 0 {
        return doc, errors.New("Empty document")
    }
    doc = m.NewDocument()

    var mark int
    %%{
        action emit_eof {
            if doc.Empty() {
                node = m.NewParagraph(data[:p])
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
            if !node.Empty() && node.Type != m.Doc {
                nodes = append(nodes, node)
                //log.Printf("appending node: %s\n", node)
                node = m.NewEmptyNode()
            }
            log.Printf("emit_add_block(%d)", p)
        }

        action emit_add_thematic_break {
            log.Printf("emit_add_thematic_break(%d)", p)
        }

        action emit_add_atx_heading {
            log.Printf("emit_add_atx_heading(%d)", p)
        }

#        single_line_doc = line_char* (eol | eop)? %eof(emit_eof);
#        document = (block %emit_add_block %mark)*;

        main := |*
#            thematic_break => emit_add_thematic_break;
#            atx_heading => emit_add_atx_heading;
#            single_line_doc => emit_add_paragraph;
            text_paragraph => emit_add_paragraph;
        *|;

#        main := document %eof(emit_eof);

        write init;
        write exec;
    }%%

    if false {
        log.Printf("last node %s", node)
        log.Printf("nodes %s", nodes)
        log.Printf("doc %s", doc)
        log.Printf("mark:%d, te:%d, act:%d", mark, te, act)
    }
    return doc, nil
}
