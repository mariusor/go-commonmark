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
    if end_of_par == 0 {
        end_of_par = p
    }
    node = m.NewParagraph(data[mark:end_of_par])
    log.Printf("par(%d): %s", end_of_par, node)
}

line = line_char+ eol;
single_line_block = line_char+ eol? >emit_add_paragraph;
text_paragraph = line_char+ eop >emit_add_paragraph;
container_block = (single_line_block | text_paragraph);

leaf_block = (thematic_break | atx_heading);
block = (container_block | leaf_block);

}%%
