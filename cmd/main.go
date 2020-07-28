// +build !unit

package main

import (
	"flag"
	"log"
	"os"
)

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
