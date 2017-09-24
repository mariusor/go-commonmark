// -*-go-*-
//
// Commonmark character level definitions
// Copyright (c) 2017 Marius Orcsik <marius@habarnam.ro>
// MIT License
// 

%%{

machine character_definitions;

action non_printable_ascii
{
    //log.Printf("np: %s\n", data[p]);
}

action two_byte_utf8_sequence
{
    //log.Printf("2b %s\n", data[p-2:p]);
}

action three_byte_utf8_sequence
{
    //log.Printf("3b %s\n", data[p-3:p]);
}

action four_byte_utf8_sequence
{
    //log.Printf("4b %s\n", data[p-4:p]);
}

action two_byte_utf8_space {
    //log.Printf("2bsp %s\n", data[p-2:p]);
}
action three_byte_utf8_space {
    //log.Printf("3bsp %s\n", data[p-3:p]);
}

action replace_insecure_char 
{
    log.Printf("insecurepos %d", p)
    data = arr_splice(data, []byte{0xef, 0xbf, 0xbd}, p)
    // readjusting the pointers, as we just resized the data buffer
    eof = len(data)
    pe = eof
}

action mark {
    log.Printf("mark(%d)", p)
    mark = p
}

action emit_add_line {
//    if !node.Empty() {
//        node.Children = append(node.Children, NewInlineText(data[mark:p]))
//    }
//    if node.Empty() {
//        node = NewParagraph(data[mark:p])
//    }
    log.Printf("emit_add_line(%d)", p)
}

action print_char {
    log.Printf("pos:%d char: %s", p, string(data[p]))
}

replacement = 0xef 0xbf 0xbd;
insecure = 0x00 >replace_insecure_char;

# http://spec.commonmark.org/0.27/#ascii-punctuation-character
# ! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~
punctuation = (0x21..0x2f | 0x3a..0x40 | 0x5b..0x60 | 0x7b..0x7e);

# all the printable ASCII characters (0x20 to 0x7e) excluding those explicitly covered elsewhere
ascii_char = (0x21..0x7e) -- punctuation;

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
utf8_char = (0x01..0x1f | 0x7f)                             %non_printable_ascii        |
            (0xc2..0xdf 0x80..0xbf)                         %two_byte_utf8_sequence     |
            (0xe0..0xef 0x80..0xbf 0x80..0xbf)              %three_byte_utf8_sequence   |
            (0xf0..0xf4 0x80..0xbf 0x80..0xbf 0x80..0xbf)   %four_byte_utf8_sequence;

# LF and CR characters
eol_char = ((0x0d? 0x0a) | 0x0d);

# UTF-8 white space characters
utf8_space = (0xc2 0xa0)               %two_byte_utf8_space    | # no-break-space 
             (0xe1 0x9a 0x80                                   | # ogham space
             (0xe2 0x80 (0x80..0x8a | 0xaf | 0x9f))            | # en/em quad, en/em space, three, four and six per em space, figure space, punctuation space, 
                                                                 #   thin space, hair space, narrow no-break space, medium mathematical space
             (0xe3 0x80 0x80))         %three_byte_utf8_space;   # ideographic space

# Space, tab and utf8 space characters -> inline space
i_space = 0x20 | 0x09 | utf8_space;

ws = i_space | eol_char;

character = ascii_char | utf8_char;
line_char = (i_space | character | insecure | punctuation);# %print_char;

eol = (eol_char{1}) %emit_add_line;

# eol terminated line
#line = line_char* eol;
}%%
