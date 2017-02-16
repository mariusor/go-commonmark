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

func Parse (data []byte) (bool, error) {

    links := parse(data)

    if links == nil {
        return false, errors.New(":")
    }
    return true, nil
}

%% machine parser;
%% include character_definitions "characters.rl";
%% write data;

func parse(data []byte) []byte{
    cs, p, pe := 0, 0, len(data)
    eof := len(data)

    if pe == 0 {
        return nil
    }

    fmt.Printf("Incoming str: %s - len %d\n", data, len(data))

    var header_level int;
    var mark int;

    %%{

        action emit_thematic_break 
        {
            fmt.Printf("thematic break\n");
        }

        action emit_header_start
        {
            fmt.Printf("h start\n")
        }

        action emit_header_level_start
        {
            mark = p;
        }

        action emit_header_level_end
        {
            header_level = p - mark;
            mark = p;

        }

        action emit_header_end
        {
            fmt.Printf("header end - level %d, cont %s\n", header_level, data[mark:p]);
        }


        tematic_break_char = ('*' | '-' | '_');
        thematic_break_underscore = (' '{1,3} ('_' | sp){3,} eol) %emit_thematic_break;
        thematic_break_star = (' '{1,3} ('*' | sp){3,} eol) %emit_thematic_break;
        thematic_break_minus = (' '{1,3} ('-' | sp){3,} eol) %emit_thematic_break;

        thematic_break = thematic_break_underscore | thematic_break_star | thematic_break_minus;
 

        header = (('#'{1,6}) >emit_header_level_start %emit_header_level_end ) %emit_header_start char*;
        atx_headings = (ws{0,3} header eol) %emit_header_end;

        block = thematic_break* | atx_headings* | line*;

        main := block*;
 
        write init;
        write exec;
    }%%

    return data 
}
