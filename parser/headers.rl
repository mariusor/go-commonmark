
// -*-go-*-
//
// Commonmark parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
// MIT License
// 

%%{

machine headers;


action emit_header_start
{
}

action emit_header_level
{
    header_level++;
}

action emit_header_level_end
{
    mark = p
}

action emit_header_end
{
    node = NewHeader(header_level, data[mark:p])
}

header = ('#'{1,6} @emit_header_level %emit_header_level_end) (sp+)%mark (char | ws)*;
headers = (ws{0,3} header) %emit_header_end eol;

#write data nofinal;
}%%
