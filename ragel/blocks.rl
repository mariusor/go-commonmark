// -*-go-*-
//
// Commonmark block definitions
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
//

%%{
machine blocks;

include chars "characters.rl";
include thematic_breaks "thematic_breaks.rl";
include headings "headings.rl";

action emit_add_paragraph {
    node = m.NewParagraph(data[mark:p])
    log.Printf("par(%d): %s", p, node)
    //mark = p
}

text_paragraph =  line_char+ (eop | eol) %emit_add_paragraph;

leaf_block = thematic_break | atx_heading;

container_block = text_paragraph;

block = leaf_block | container_block;

}%%
