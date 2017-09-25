// -*-go-*-
//
// Commonmark thematic breaks
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
//

%%{

machine thematic_breaks;

action save_break_symbol {
    thematic_break_symbol = data[p]
}

action emit_thematic_break 
{
    node = m.NewThematicBreak(thematic_break_symbol)
    log.Printf("th-br(%d):%s", p, node);
}

thematic_break_underscore = (i_space{1,3} ('_' | i_space){3,} >save_break_symbol);
thematic_break_star = (i_space{1,3} ('*' | i_space){3,} >save_break_symbol);
thematic_break_minus = (i_space{1,3} ('-' | i_space){3,} >save_break_symbol);

thematic_break = ((thematic_break_underscore | thematic_break_star | thematic_break_minus) (eop | eol)? %emit_thematic_break);

}%%
