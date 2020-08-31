package main

import (
	"bytes"
	"testing"
)

func testRun(opt options) (string, exitCode) {
	var output bytes.Buffer
	exitCode := run(&output, opt)
	return output.String(), exitCode
}

func assertCode(t *testing.T, opt options, want exitCode) {
	_, got := testRun(opt)

	if got != want {
		t.Errorf("Input: %v Want: %d Got: %d", opt, want, got)
	}
}

func assert(t *testing.T, opt options, wantCode exitCode, want string) {
	got, gotCode := testRun(opt)

	if gotCode != wantCode || got != want {
		t.Errorf("Input: %v Want: %d '%s' Got: %d '%s'", opt, wantCode, want, gotCode, got)
	}
}

func TestPass(t *testing.T) {
	assertCode(t, options{
		directory: "../../test/data/pass",
		format:    "",
	}, successCode)
}

func TestNoOwners(t *testing.T) {
	assert(t, options{
		directory: "../../test/data/no_owners",
		format:    "",
	}, errorCode, `CODEOWNERS 1 ::Error:: No owners specified [NoOwner]
`)
}

func TestCustomFormat(t *testing.T) {
	assert(t, options{
		directory: "../../test/data/noowners",
		format:    "test",
	}, errorCode, `test
`)
}

func TestInvalidFormatParse(t *testing.T) {
	assertCode(t, options{
		directory: "../../test/data/noowners",
		format:    "  {{template \"one ",
	}, unexpectedErrorCode)
}

func TestInvalidFormatExec(t *testing.T) {
	assertCode(t, options{
		directory: "../../test/data/noowners",
		format:    "  {{template \"one\"}}  ",
	}, unexpectedErrorCode)
}

func TestInvalidDirectory(t *testing.T) {
	assert(t, options{
		directory: "'",
		format:    "",
	}, errorCode, `CODEOWNERS 0 ::Error:: No CODEOWNERS file found [NoCodeowners]
`)
}

func TestMultipleCodeOwners(t *testing.T) {
	assert(t, options{
		directory: "../../test/data/multiple_codeowners",
		format:    "",
	}, warningCode, `CODEOWNERS 0 ::Warning:: Multiple CODEOWNERS files found (CODEOWNERS, docs/CODEOWNERS) [MultipleCodeowners]
`)
}
