
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

action emit_end_of_heading {
    end_of_heading = p
    log.Printf("end heading %d(%d)", heading_level, p)
}

action emit_heading_end
{
    if end_of_heading == 0 {
        end_of_heading = p
    }
    node = m.NewHeading(heading_level, bytes.Trim(data[mark:end_of_heading], " \n\r"))
    log.Printf("h%d(%d): %s", heading_level, p, node)
    end_of_heading = 0
}

heading_symbol = 0x23; # '#'
heading_level = (heading_symbol{1,6} @emit_heading_level);
heading_end = heading_symbol{1,6};
heading_char = line_char | punctuation;

heading = heading_level i_space+ %emit_heading_level_end heading_char+ (i_space+ heading_end+ >emit_end_of_heading)?;
atx_heading = ((i_space{0,3} heading) (eop | eol) %emit_heading_end);

}%%
