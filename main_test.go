package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDuplicateFileWithAllOptionsAndAddOptional(t *testing.T) {
	m := &Modifier{
		fieldToAddPostPreFix: make(map[string]bool),
		goPackageName:        "github.com/davidae/faker;faker",
		packageName:          "faker",
		prefix:               "Pre",
		postfix:              "Post",
		addOptional:          true,
	}

	input, err := os.Open("test/example.proto")
	if err != nil {
		t.Fatal("failed to open proto test file")
	}

	expectedOut, err := ioutil.ReadFile("test/add_optional.proto")
	if err != nil {
		t.Fatal("failed to open proto expected output test file")
	}

	out := duplicateFile(m, input)
	if out != string(expectedOut) {
		t.Fatalf("expected output isn't equal expected output, expected:\n%s\nactual:\n%s\n", string(expectedOut), out)
	}
}
func TestDuplicateFileWithAllOptionsAndRemoveOptional(t *testing.T) {
	m := &Modifier{
		fieldToAddPostPreFix: make(map[string]bool),
		prefix:               "Pre",
		removeOptional:       true,
	}

	input, err := os.Open("test/example.proto")
	if err != nil {
		t.Fatal("failed to open proto test file")
	}

	expectedOut, err := ioutil.ReadFile("test/remove_optional.proto")
	if err != nil {
		t.Fatal("failed to open proto expected output test file")
	}

	out := duplicateFile(m, input)
	if out != string(expectedOut) {
		t.Fatalf("expected output isn't equal expected output, expected:\n%s\nactual:\n%s\n", string(expectedOut), out)
	}
}
