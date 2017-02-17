package main

import (
    "./parser"
    "testing"
    "os"
    "io"
    "bytes"
)

type testPair struct {
    text        string
    expected    bool
    doc         parser.Document
}

var emptyDoc parser.Document = parser.Document{}

var someTests = []testPair{
//    {
//        {0x00},
//        false,
//        emptyDoc,
//    },
    {
        "[ana](httpslittrme)",
        true,
        emptyDoc,
    },
    {
        "[ana](https://littr.me)",
        true,
        emptyDoc,
    },
    {
        "",
        false,
        emptyDoc,
    },
//    { "some text before [test 123](https://littr.me)", true },
//    { "[test 123](https://littr.me) some text after", true },
//    { "some text before [test 123](https://littr.me) some text after", true },
    {
        "𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș\n", 
        true ,
        parser.Document{
            Children: []parser.Node{
                parser.Node{
                    Type: parser.Par,
                    Content: []byte("𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș"),
                },
            },
        },
    },
//    { " ---\n", true },
//    { "  ***\n", true },
//    { "  * * * *\n", true },
//    { "   ___\r", true },
//    { "   _*-*__\r", true },
    {
        "# ana are mere\n", 
        true,
        parser.Document{
            Children: []parser.Node{
                parser.Node{
                    Type: parser.H1,
                    Content: []byte("ana are mere"),
                },
            },
        },
    },
    {
        "## ana are mere\n", 
        true,
        parser.Document{
            Children: []parser.Node{
                parser.Node{
                    Type: parser.H2,
                    Content: []byte("ana are mere"),
                },
            },
        },
    },

    {
        "### ana are mere\n",
        true,
        parser.Document{
            Children: []parser.Node{
                parser.Node{
                    Type: parser.H3,
                    Content: []byte("ana are mere"),
                },
            },
        },
    },
    {
        "#### ana are mere\n",
        true,
        parser.Document{
            Children: []parser.Node{
                parser.Node{
                    Type: parser.H4,
                    Content: []byte("ana are mere"),
                },
            },
        },
    },
    {
        "##### ana-are-mere\n",
        true,
        parser.Document{
            Children: []parser.Node{
                parser.Node{
                    Type: parser.H5,
                    Content: []byte("ana-are-mere"),
                },
            },
        },
    },
    {
        "###### ana-are-mere\n",
        true,
        parser.Document{
            Children: []parser.Node{
                parser.Node{
                    Type: parser.H6,
                    Content: []byte("ana-are-mere"),
                },
            },
        },
    },
}

func TestParse (t *testing.T) {
    for _, curTest := range someTests {
        doc, err  := parser.Parse([]byte(curTest.text))
        t.Logf("Testing '%s'.", curTest.text)
        if err != nil && curTest.expected {
            t.Errorf("\tParse result invalid, expected '%t, got %v'\n", curTest.expected)
        }
        if curTest.doc.Equal(doc) {
            t.Logf("\t%s\n", doc.Children)
        }
        testChildren := curTest.doc.Children
        children := doc.Children
        if len(testChildren) != len(children) {
            t.Errorf("\tParse result invalid, children length expected %d != %d\n", len(testChildren), len(children))
        }
        if len(testChildren) > 0 {
            testNode := testChildren[0]
            node := children[0]
            if testNode.Type != node.Type {
                t.Errorf("\tParse result invalid, node type expected %v != %v\n", testNode.Type, node.Type)
            }
            if !bytes.Equal(testNode.Content, node.Content) {
                t.Errorf("\tParse result invalid, node content expected %s != %s\n", testNode.Content, node.Content)
                t.Errorf("\tParse result invalid, node content expected %v != %v\n", testNode.Content, node.Content)
            }
            //t.Logf("\t%s\n", testNode.Content)
        }
    }
}

func TestParseReadme (t *testing.T) {
    f, _ := os.Open("README.md")

    data := make([]byte, 512)
    io.ReadFull(f, data)

    if _, err := parser.Parse(data); err != nil {
        t.Errorf("\tParse invalid\n")
    }
}

func TestMain(m *testing.M) {
    os.Exit(m.Run())
}
