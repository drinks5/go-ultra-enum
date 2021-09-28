package main

import (
	"flag"

	"github.com/drinks5/go-ultra-enum/parser"
)

var (
	file   = flag.String("file", "", "The file(s) to generate enums.  Use more than one flag for more files.")
	output = flag.String("output", "", "output file name; default srcdir/enumer.go")
)

func main() {
	p := parser.Parser{}
	flag.Parse()

	p.Packages(*file)
	p.Render(parser.Header, p.Pkg)
	p.Generate()
	p.Sink(*output)
}
