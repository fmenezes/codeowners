package main

import (
	"flag"
	"io"
	"log"
	"os"
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

	format := "{{range .}}{{ .LineNo }} ::{{ .Severity }}:: {{ .Message }} [{{ .CheckName }}]\n{{end}}"
	if len(opt.format) > 0 {
		format = opt.format
	}
	tpl, err := template.New("main").Parse(format)
	if err != nil {
		return err
	}

	checkers := codeowners.AvailableCheckers()

	checks, err := codeowners.Check(dir, checkers...)
	if err != nil {
		return err
	}

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

func main() {
	dir := flag.String("d", ".", "Specifies the directory you want to use to lint the CODEOWNERS file")
	format := flag.String("f", "", "Specifies the format you want to return lint results")
	flag.Parse()

	opt := options{
		directory: *dir,
		format:    *format,
	}

	if err := run(os.Stdout, opt); err != nil {
		log.Fatal(err)
	}
}
