
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
    node = NewHeading(heading_level, data[mark:p])
}

heading_level = ('#'{1,6} @emit_heading_level %emit_heading_level_end);
heading_char = (line_char | asciipunct);

heading = heading_level sp* (heading_char+ >mark);
atx_heading = (ws{0,3} heading) eol >emit_heading_end eol?;

#write data nofinal;
}%%
