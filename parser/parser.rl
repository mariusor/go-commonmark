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
%% include thematic_breaks "thematic_breaks.rl";
%% include headers "headers.rl";
%% write data;

func parse(data []byte) []byte{
    cs, p, pe := 0, 0, len(data)
    eof := len(data)

    if pe == 0 {
        return nil
    }

    // fmt.Printf("Incoming str: %s - len %d\n", data, len(data))

    var header_level int
    var mark int
    var header_cont []byte

    %%{

        block = thematic_break* | atx_headings* | line*;

        main := block*;
 
        write init;
        write exec;
    }%%

    return data 
}
