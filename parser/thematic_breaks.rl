// -*-go-*-
//
// Commonmark parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
// MIT License
// 

%%{

machine thematic_breaks;


action emit_thematic_break 
{
    node = NewThematicBreak('-')
}

thematic_break_underscore = (' '{1,3} ('_' | sp){3,} eol) %emit_thematic_break;
thematic_break_star = (' '{1,3} ('*' | sp){3,} eol) %emit_thematic_break;
thematic_break_minus = (' '{1,3} ('-' | sp){3,} eol) %emit_thematic_break;

thematic_break = thematic_break_underscore | thematic_break_star | thematic_break_minus;

#write data nofinal;
}%%
