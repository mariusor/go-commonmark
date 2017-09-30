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

func load_files(ext string) ([]string, error) {
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

func get_file_contents(path string) []byte {
	f, _ := os.Open(path)

	data := make([]byte, 512)
	io.ReadFull(f, data)
	data = bytes.Trim(data, "\x00")

	return data
}

func TestWithFiles(t *testing.T) {
	var tests []string
	var err error

	tests, err = load_files(".md")
	var f = t.Errorf
	if stopOnFailure {
		f = t.Fatalf
	}
	if len(tests) == 0 {
		t.Skip("No tests found")
	}
	for _, path := range tests {
		var doc m.Document
		var t_doc testDoc
		if path == "../../tests/README.md" {
			continue
		}
		data := get_file_contents(path)
		t.Logf("Testing: %s", path)
		res_path := fmt.Sprintf("%s.json", path[:len(path)-3])
		json.Unmarshal(get_file_contents(res_path), &t_doc)

		res_doc := t_doc.Document()
		doc, err = Parse(data)

		if err != nil {
			f("%s", err)
		}
		if !reflect.DeepEqual(res_doc, doc) {
			f("\n%#v\n%#v\n%s\n%s", doc, res_doc, doc, res_doc)
		}
	}
}
