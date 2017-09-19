
// -*-go-*-
//
// Commonmark headings definitions
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
// 

%%{

machine headings;


action emit_heading_start
{
}

action emit_heading_level
{
    heading_level++;
}

action emit_heading_level_end
{
    log.Printf("hle: %d\n", p)
    mark = p
}

action emit_heading_end
{
    log.Printf("he: %d:%d\n", mark, p)
    node = NewHeading(heading_level, data[mark:p])
}

heading_symbol = 0x23;
heading_level = (heading_symbol{1,6} @emit_heading_level);
heading_end = i_space? heading_symbol*;
heading_char = (line_char | punctuation);

heading = heading_level i_space+ %emit_heading_level_end (heading_char+) heading_end;
atx_heading = (i_space{0,3} heading) eol >emit_heading_end;

}%%
