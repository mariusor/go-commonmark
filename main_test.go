package main

import (
    "./parser"
    "testing"
    "os"
    "io"
)

type testPair struct {
    text        string
    expected    bool
}

var someTests = []testPair{
    { "[ana](httpslittrme)", true },
    { "[ana](https://littr.me)", true },
    { "", false },
    { "some text before [test 123](https://littr.me)", true },
    { "[test 123](https://littr.me) some text after", true },
    { "some text before [test 123](https://littr.me) some text after", true },
    { "𐍈ᏚᎢᎵᎬᎢᎬᏒăîțș", true },
    { " ---\n", true },
    { "  ***\n", true },
    { "  * * * *\n", true },
    { "   ___\r", true },
    { "   _*-*__\r", true },
    { "# ana-are-mere\n", true },
    { "## ana-are-mere\n", true },
    { "### ana-are-mere\n", true },
    { "#### ana-are-mere\n", true },
    { "##### ana-are-mere\n", true },
    { "###### ana-are-mere\n", true },
}

func TestParse (t *testing.T) {
    for _, curTest := range someTests {
        oops, _  := parser.Parse([]byte(curTest.text))
        t.Logf("Testing '%s'.", curTest.text)
        if oops != curTest.expected {
            t.Errorf("\tParse result invalid, expected '%t'\n", curTest.expected)
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
