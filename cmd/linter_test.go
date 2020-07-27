package main

import (
	"bytes"
	"testing"
)

func testRun(opt options) (string, error) {
	var output bytes.Buffer
	err := run(&output, opt)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func TestSimple(t *testing.T) {
	opt := options{
		directory: "../test/data/simple",
		format:    "",
	}

	got, err := testRun(opt)
	if err != nil {
		t.Error(err)
	}

	want := `Everything ok ;)
`
	if got != want {
		t.Errorf("Input: %v Want: '%s' Got: '%s'", opt, want, got)
	}
}

func TestNoOwners(t *testing.T) {
	opt := options{
		directory: "../test/data/noowners",
		format:    "",
	}

	got, err := testRun(opt)
	if err != nil {
		t.Error(err)
	}

	want := `1 ::Error:: No owners specified [NoOwner]
`
	if got != want {
		t.Errorf("Input: %v Want: '%s' Got: '%s'", opt, want, got)
	}
}

func TestCustomFormat(t *testing.T) {
	opt := options{
		directory: "../test/data/noowners",
		format:    "test",
	}

	got, err := testRun(opt)
	if err != nil {
		t.Error(err)
	}

	want := `test`
	if got != want {
		t.Errorf("Input: %v Want: '%s' Got: '%s'", opt, want, got)
	}
}

func TestInvalidFormat(t *testing.T) {
	opt := options{
		directory: "../test/data/noowners",
		format:    "  {{template \"one\"}}  ",
	}

	_, err := testRun(opt)
	if err == nil {
		t.Errorf("Should have errored")
	}
}

func TestInvalidDirectory(t *testing.T) {
	opt := options{
		directory: "../test/da\\ata",
		format:    "",
	}

	got, err := testRun(opt)
	if err != nil {
		t.Error(err)
	}

	want := `0 ::Error:: No CODEOWNERS file found [NoCodeowners]
`
	if got != want {
		t.Errorf("Input: %v Want: '%s' Got: '%s'", opt, want, got)
	}
}
