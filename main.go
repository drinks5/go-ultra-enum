package main

import (
	"flag"

	"github.com/drinks5/go-ultra-enum/parser"
)

var (
	output = flag.String("output", "", "output file name; default srcdir/enumer.go")
)

func main() {
	p := parser.Parser{}
	args := flag.Args()
	p.Packages(args)
	p.Render(parser.Header, p.Pkg)
	p.Generate()
	p.Sink(*output)
}
