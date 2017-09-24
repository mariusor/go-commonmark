package parser

import (
	"testing"
	"os"
	"fmt"
	"encoding/json"
	m "markdown"
	"path/filepath"
	"errors"
	"io"
	"bytes"
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

func testWithFiles(t *testing.T) {
	var tests []string
	var err error

	tests, err = load_files(".md")
	f := t.Errorf
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
		j_res_doc, _ := json.Marshal(res_doc)
		j_doc, _ := json.Marshal(doc)
		if string(j_res_doc) != string(j_doc) {
			f("\n%s_________________\n%s", doc, res_doc)
		}
	}
}
