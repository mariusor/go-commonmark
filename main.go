package main

import (
	"./parser"
	"bytes"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
)

const (
	version = "0.0.1"
)

var (
	fileName = kingpin.Arg("file path", "The path of the file to display.").Required().String()
)

var exitWithError = func(e error) {
	fmt.Printf("error: %s\n", e)
	os.Exit(1)
	return
}

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version(version)
	kingpin.Parse()

	var status int = 0
	f, err := os.Open(*fileName)
	if err != nil {
		exitWithError(err)
	}

	data := make([]byte, 512)
	io.ReadFull(f, data)
	data = bytes.Trim(data, "\x00")

	doc, err := parser.Parse(data)
	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("%s\n", doc.String())

	os.Exit(status)
	return
}
