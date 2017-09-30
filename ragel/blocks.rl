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

single_line_doc = ((line_char | punctuation)** (eol)?);
text_paragraph = ((line_char | punctuation)** eop);

leaf_block = thematic_break | atx_heading;
container_block = text_paragraph;
block = leaf_block | container_block;

}%%
