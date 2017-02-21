
// -*-go-*-
//
// Commonmark parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
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
    mark = p
}

action emit_heading_end
{
    if mark > 0 {
        node = NewHeading(heading_level, data[mark:p])
        mark = -1
    }
}

heading_symbol = 0x23;
heading_level = (heading_symbol{1,6} @emit_heading_level %emit_heading_level_end);
heading_char = (line_char | punctuation);

heading = heading_level i_space+ (heading_char+ >mark);
atx_heading = (i_space{0,3} heading) eol >emit_heading_end eol?;

#write data nofinal;
}%%
