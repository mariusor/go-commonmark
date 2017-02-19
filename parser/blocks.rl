// -*-go-*-
//
// Commonmark link parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
// MIT License
// 

%%{
machine blocks;

action emit_paragraph {
    node.children = append(node.children, NewParagraph(data[mark:p]) 
}

include thematic_breaks "thematic_breaks.rl";
include headings "headings.rl";

paragraph = line* eol;

leaf_block = thematic_break | atx_heading;

container_block = paragraph;

block = (leaf_block | container_block);

}%%
