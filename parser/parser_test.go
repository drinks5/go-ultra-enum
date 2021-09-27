package parser

import (
	"testing"
)

//go:generate ../go-ultra-enum -f=$GOFILE --marshal --lower --ptr
type ColorEnum struct {
	Red       string `enum:"RED"`
	LightBlue string `enum:"LIGHT_BLUE"`
}

func TestParser_packages(t *testing.T) {
	p := Parser{}
	p.Packages([]string{})
	p.Generate()
	p.Sink("./")
}
