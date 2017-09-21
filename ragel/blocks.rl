// -*-go-*-
//
// Commonmark block definitions
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
//

%%{
machine blocks;

action emit_paragraph {
    if mark != p {
        par := NewParagraph(data[mark:p])
        if !node.Empty() {
            log.Printf("%s", node)
            node.Children = append(node.Children, par)
        } else {
            nodes = append(nodes, par)
        }
    }
}

include chars "characters.rl";
include thematic_breaks "thematic_breaks.rl";
include headings "headings.rl";

text_paragraph = line_char* eop >emit_paragraph;

leaf_block = thematic_break | atx_heading;

container_block = text_paragraph;

block = leaf_block | container_block;

}%%
