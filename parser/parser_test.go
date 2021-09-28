package parser

import (
	"testing"
)

func TestParser_packages(t *testing.T) {
	p := Parser{}
	p.Packages("")
	p.Generate()
	p.Sink("")
}
