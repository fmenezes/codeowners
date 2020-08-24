package main

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/fmenezes/codeowners"
	_ "github.com/fmenezes/codeowners/checkers"
)

type options struct {
	directory string
	format    string
	token     string
	tokenType string
}

type exitCode int

const (
	successCode exitCode = iota
	warningCode
	errorCode
	unexpectedErrorCode
)

func run(wr io.Writer, opt options) exitCode {
	dir, err := filepath.Abs(opt.directory)
	if err != nil {
		fmt.Fprintf(wr, "Unexpected error when parsing directory: %v", err)
		return unexpectedErrorCode
	}

	format := "{{ .Position }} ::{{ .Severity }}:: {{ .Message }} [{{ .CheckName }}]"
	if len(opt.format) > 0 {
		format = opt.format
	}
	format = fmt.Sprintf("%s\n", format)
	tpl, err := template.New("main").Parse(format)
	if err != nil {
		fmt.Fprintf(wr, "Unexpected error when parsing format: %v", err)
		return unexpectedErrorCode
	}

	checkers := codeowners.AvailableCheckers()

	checks, _ := codeowners.Check(codeowners.CheckOptions{
		Directory:       dir,
		Checkers:        checkers,
		GithubToken:     opt.token,
		GithubTokenType: opt.tokenType,
	})

	code := successCode
	for _, check := range checks {
		err = tpl.Execute(wr, check)
		if err != nil {
			fmt.Fprintf(wr, "Unexpected error when writing results: %v", err)
			return unexpectedErrorCode
		}
		switch check.Severity {
		case codeowners.Error:
			if code < errorCode {
				code = errorCode
			}
		case codeowners.Warning:
			if code < warningCode {
				code = warningCode
			}
		}
	}

	return code
}
