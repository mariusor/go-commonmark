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
    // need to find a good way to insert two bytes in the place of the null char
    // this requires in place array resize :D
    //data[p-1] = 0xff
    //data[p] = 0xfd
    data[p] = 0x3f;
}

action mark {
    //fmt.Printf("cur: %d - %s\n", p, cs)
    mark = p
}

action emit_new_line {
    //fmt.Printf("nl: %d\n", p)
    if mark > 0 {
        node = NewParagraph(data[mark:p]) 
        mark = -1
    }
}

action two_byte_utf8_space {
//    fmt.Printf("2bsp %s\n", data[p-2:p]);
}
action three_byte_utf8_space {
//    fmt.Printf("3bsp %s\n", data[p-3:p]);
}

insecure = 0x00 %replace_insecure_char;

# http://spec.commonmark.org/0.27/#ascii-punctuation-character
# ! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~
asciipunct = (0x21..0x2f | 0x3a..0x40 | 0x5b..0x60 | 0x7b..0x7e);

# all the printable ASCII characters (0x20 to 0x7e) excluding those explicitly covered elsewhere
asciic = (0x21..0x7e) -- asciipunct;

# @see: https://git.wincent.com/wikitext.git/blob/4bb2e23eebaf25c6f1dddb721f074f69375d222a:/ext/wikitext/wikitext_ragel.rl
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
eol = (0x0a | 0x0d);

# UTF-8 white space characters
utf8sp = (0xc2 0xa0)                              %two_byte_utf8_space       | # no-break-space 
         ( 0xe1 0x9a 0x80                                                    | # ogham space
           0xe2 0x80 (0x80..0x8a | 0xaf | 0x9f)                              | # en and em quad, en and em space, three, four and six per em space, figure space, punctuation space
                                                                               #    thin space, hair space, narrow no-break space, medium mathematical space
           0xe3 0x80 0x80)                        %three_byte_utf8_space;      # ideographic space

# Space and tab characters
sp = 0x20 | 0x09 | utf8sp;

#wsp = sp | eol;

char = asciic | utf8c;

line_char = (sp | char | insecure);
line = line_char* >mark eol;
#line = line_char* >mark eol %emit_new_line;

#write data nofinal;
}%%
