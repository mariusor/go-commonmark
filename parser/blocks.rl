// -*-go-*-
//
// Commonmark link parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
// MIT License
// 

%%{
machine blocks;

include thematic_breaks "thematic_breaks.rl";
include headers "headings.rl";

paragraph = (line eol)* eol | line;

leaf_block = thematic_break | atx_heading;

container_block = paragraph;

block = leaf_block | container_block;

}%%
