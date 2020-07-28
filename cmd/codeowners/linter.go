package main

import (
	"io"
	"path/filepath"
	"text/template"

	"github.com/fmenezes/codeowners"
	_ "github.com/fmenezes/codeowners/checkers"
)

type options struct {
	directory string
	format    string
}

func run(wr io.Writer, opt options) error {
	dir, err := filepath.Abs(opt.directory)
	if err != nil {
		return err
	}

	format := "{{range .}}{{ .Position }} ::{{ .Severity }}:: {{ .Message }} [{{ .CheckName }}]\n{{end}}"
	if len(opt.format) > 0 {
		format = opt.format
	}
	tpl, err := template.New("main").Parse(format)
	if err != nil {
		return err
	}

	checkers := codeowners.AvailableCheckers()

	checks, _ := codeowners.Check(dir, checkers...)

	if len(checks) > 0 {
		err = tpl.Execute(wr, checks)
		if err != nil {
			return err
		}
	} else {
		wr.Write([]byte("Everything ok ;)\n"))
	}

	return nil
}
