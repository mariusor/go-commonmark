
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
    //log.Printf("hle: %d", p)
    mark = p
}

action emit_heading_end
{
    node = m.NewHeading(heading_level, bytes.Trim(data[mark:p], "\n\r"))
    log.Printf("h%d(%d): %s", heading_level, p, node)
}

heading_symbol = 0x23; # '#'
heading_level = (heading_symbol{1,6} @emit_heading_level);
heading_end = i_space? heading_symbol*;
heading_char = (line_char | punctuation | ^eol);

heading = heading_level i_space+ %emit_heading_level_end (heading_char+) heading_end;
atx_heading = ((i_space{0,3} heading) (eop | eol)? %emit_heading_end);

}%%
