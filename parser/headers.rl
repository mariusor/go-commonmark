
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
    header_level += 1;
}

action emit_header_level_end
{
    mark = p
    header_level = 0

}

action emit_header_end
{
    header_cont = data[mark:p]
    fmt.Printf("header end - level %d, cont %s\n", header_level, header_cont);
}

header_char = char -- '#';
header = ('#'{1,6} @emit_header_level %emit_header_level_end) sp+ %emit_header_start char* ;
atx_headings = (ws{0,3} header) %emit_header_end eol;

#write data nofinal;
}%%
