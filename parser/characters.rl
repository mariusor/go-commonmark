// -*-go-*-
//
// Commonmark parser
// Copyright (c) 2017 Marius Orcisk <marius@habarnam.ro>
// MIT License
// 

%%{

machine character_definitions;

action non_printable_ascii
{
    //fmt.Printf("np: %s\n", data[p]);
}

action two_byte_utf8_sequence
{
    //fmt.Printf("2b %s\n", data[p-2:p]);
}

action three_byte_utf8_sequence
{
    //fmt.Printf("3b %s\n", data[p-3:p]);
}

action four_byte_utf8_sequence
{
    //fmt.Printf("4b %s\n", data[p-4:p]);
}

action replace_insecure_char 
{
    // need to find a good way to insert two bytes in the place of the faulty char
    // this requires in place array resize :D
    //data[p-1] = 0xff
    //data[p] = 0xfd
    data[p] = byte('?');
}

action mark {
    mark = p
}

action emit_new_line {
    node = NewParagraph(data[mark:p]) 
}

insecure = (0x00 0x00) %replace_insecure_char;

# all the printable ASCII characters (0x20 to 0x7e) excluding those explicitly covered elsewhere:
# skip space (0x20), quote (0x22), ampersand (0x26), less than (0x3c), greater than (0x3e),
# left bracket 0x5b, right bracket 0x5d, backtick (0x60), and vertical bar (0x7c)
asciic = (0x21 | 0x23..0x25 | 0x27..0x3b | 0x3d | 0x3f..0x5a | 0x5c | 0x5e..0x5f | 0x61..0x7b | 0x7e);

# here is where we handle the UTF-8 and everything else
#
#     one_byte_sequence   = byte begins with zero;
#     two_byte_sequence   = first byte begins with 110 (0xc0..0xdf), next with 10 (0x80..9xbf);
#     three_byte_sequence = first byte begins with 1110 (0xe0..0xef), next two with 10 (0x80..9xbf);
#     four_byte_sequence  = first byte begins with 11110 (0xf0..0xf7), next three with 10 (0x80..9xbf);
#
#     within the ranges specified, we also exclude these illegal sequences:
#       1100000x (c0 c1)    overlong encoding, lead byte of 2 byte seq but code point <= 127
#       11110101 (f5)       restricted by RFC 3629 lead byte of 4-byte sequence for codepoint above 10ffff
#       1111011x (f6, f7)   restricted by RFC 3629 lead byte of 4-byte sequence for codepoint above 10ffff
utf8c = (0x01..0x1f | 0x7f)                             %non_printable_ascii        |
        (0xc2..0xdf 0x80..0xbf)                         %two_byte_utf8_sequence     |
        (0xe0..0xef 0x80..0xbf 0x80..0xbf)              %three_byte_utf8_sequence   |
        (0xf0..0xf4 0x80..0xbf 0x80..0xbf 0x80..0xbf)   %four_byte_utf8_sequence;

# LF and CR characters
eol = (0x0d? 0x0a) | 0x0d;

# UTF-8 white space characters
utf8sp = (0xc2 0xa0)                              %two_byte_utf8_space       | # no-break-space 
         ( 0xe1 0x9a 0x80                                                    | # ogham space
           0xe2 0x80 (0x80..0x8a | 0xaf | 0x9f)                              | # en and em quad, en and em space, three, four and six per em space, figure space, punctuation space
                                                                               #    thin space, hair space, narrow no-break space, medium mathematical space
           0xe3 0x80 0x80)                        %three_byte_utf8_space;      # ideographic space

# Space and tab characters
sp = 0x20 | 0x09 | utf8sp;

ws = sp | eol;

char = asciic | utf8c;

#line = (char asciipunct)* >mark %emit_new_line eol;
line = (char)* >mark %emit_new_line eol;

#write data nofinal;
}%%
