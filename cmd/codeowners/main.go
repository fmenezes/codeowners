// +build !unit

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	dir := flag.String("d", ".", "Directory: specifies the directory you want to use to lint the CODEOWNERS file")
	format := flag.String("f", "", "Format: specifies the format you want to return lint results")
	flag.Parse()

	opt := options{
		directory: *dir,
		format:    *format,
	}
	exitCode := run(os.Stderr, opt)
	if exitCode == successCode {
		fmt.Println("Everything ok ;)")
		return
	}
	os.Exit(int(exitCode))
}
