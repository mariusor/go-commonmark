package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	m "markdown"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func loadFiles(ext string) ([]string, error) {
	var files []string

	dir := "../../tests"
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && path[len(path)-len(ext):] == ext {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not load files %s/*%s", dir, ext))
	}

	return files, nil
}

func getFileContents(path string) []byte {
	f, _ := os.Open(path)

	data := make([]byte, 512)
	io.ReadFull(f, data)
	data = bytes.Trim(data, "\x00")

	return data
}

func _TestWithFiles(t *testing.T) {
	var tests []string
	var err error

	tests, err = loadFiles(".md")
	var f = t.Errorf
	if stopOnFailure {
		f = t.Fatalf
	}
	if len(tests) == 0 {
		t.Skip("No tests found")
	}
	for _, path := range tests {
		var doc m.Document
		var tDoc testDoc
		if path == "../../tests/README.md" {
			continue
		}
		data := getFileContents(path)
		t.Logf("Testing: %s", path)
		resPath := fmt.Sprintf("%s.json", path[:len(path)-3])
		json.Unmarshal(getFileContents(resPath), &tDoc)

		resDoc := tDoc.Document()
		doc, err = Parse(data)

		if err != nil {
			f("%s", err)
		}
		if !reflect.DeepEqual(resDoc, doc) {
			f("\n%#v\n%#v\n%s\n%s", doc, resDoc, doc, resDoc)
		}
	}
}
