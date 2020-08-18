[![Go](https://github.com/fmenezes/codeowners/workflows/Go/badge.svg)](https://github.com/fmenezes/codeowners/actions?query=branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/fmenezes/codeowners)](https://goreportcard.com/report/github.com/fmenezes/codeowners)
[![Coverage](https://coveralls.io/repos/github/fmenezes/codeowners/badge.svg?branch=master)](https://coveralls.io/github/fmenezes/codeowners?branch=master)
[![Godoc](https://godoc.org/github.com/fmenezes/codeowners?status.svg)](https://godoc.org/github.com/fmenezes/codeowners)

# CODEOWNERS

CodeOwners package provides funcionality to evaluate CODEOWNERS file in Go. Also provides a CLI to lint.

## Documentation

### Package

To find package documentation follow https://godoc.org/github.com/fmenezes/codeowners

### CLI

#### Installation

Simply run `go get -u github.com/fmenezes/codeowners/cmd/codeowners`

#### Usage

Simply calling `codeowners` will kick off the cli on the current directory.

##### Options

| Option        | Default Value | Description                                                                    |
| ------------- | ------------- | ------------------------------------------------------------------------------ |
| d             | .             | Directory: specifies the directory you want to use to lint the CODEOWNERS file |
| f             |               | Format: specifies the format you want to return lint results                   |

##### Exit Codes

| Exit Code     | Description                                                      |
| ------------- | ---------------------------------------------------------------- |
| 0             | Success: no errors returned                                      |
| 1             | Warnings: linter returned a few warnings but no errors           |
| 2             | Errors: linter returned a few errors                             |
| 3             | Unexpected errors: errors that prevented the linter from running |

## Compatibility

:warning: This module is on a v0 mode and it is not ready to be used, once it reaches the v1 we will lock the API.
