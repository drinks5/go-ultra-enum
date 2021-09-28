package parser

import (
	"testing"
)

func TestParser_packages(t *testing.T) {
	p := Parser{}
	p.Packages("")
	p.Render(Header, p.Pkg)
	p.Generate()
	p.Sink("")
}
