package parser

import (
	"testing"
)

func TestParser_packages(t *testing.T) {
	p := Parser{}
	p.Packages("")
	p.Generate("")
	// r, _ := Color.Red.MarshalJson()
	// fmt.Println(string(r))
}
