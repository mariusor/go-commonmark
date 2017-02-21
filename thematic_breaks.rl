// -*-go-*-
//
// Commonmark parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
// MIT License
// 

%%{

machine thematic_breaks;

action save_break_symbol {
    thematic_break_symbol = data[p]
}

action emit_thematic_break 
{
    node = NewThematicBreak(thematic_break_symbol)
}

thematic_break_underscore = (i_space{1,3} ('_' | i_space){3,} >save_break_symbol eol) %emit_thematic_break;
thematic_break_star = (i_space{1,3} ('*' | i_space){3,} >save_break_symbol eol) %emit_thematic_break;
thematic_break_minus = (i_space{1,3} ('-' | i_space){3,} >save_break_symbol eol) %emit_thematic_break;

thematic_break = thematic_break_underscore | thematic_break_star | thematic_break_minus;

}%%
