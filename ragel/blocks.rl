// -*-go-*-
//
// Commonmark block definitions
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
//

%%{
machine blocks;

action emit_add_paragraph {
    if mark != p {
        node = NewParagraph(bytes.Trim(data[mark:p], "\n\r"))
    }
    log.Printf("emit_add_paragraph(%d)", p)
    mark = p
}

#action emit_add_paragraph {
#    if node.Empty() {
#        node = NewParagraph(data[mark:p])
#    }
#    log.Printf("emit_new_paragraph(%d)", p)
#    node = Node{}
#}

include chars "characters.rl";
include thematic_breaks "thematic_breaks.rl";
include headings "headings.rl";

#eop = ((0x0d 0x0a 0x0d 0x0a) | (0x0d 0x0d) | (0x0a 0x0a)) %emit_add_paragraph;
eop = (eol{2,}) %emit_add_paragraph;

text_paragraph =  line_char+ eop;

leaf_block = thematic_break | atx_heading;

container_block = text_paragraph*;

block = leaf_block | container_block;
#block = container_block;

}%%
