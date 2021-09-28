package parser

//go:generate go-ultra-enum
type ColorEnum struct {
	Red       string `enum:"RED"`
	LightBlue string `enum:"LIGHT_BLUE"`
}

// generate an enum with default display values. The display values are set to the field names, e.g. `On` and `Off`
type StatusEnum struct {
	On  int `enum:"1"`
	Off int `enum:"0"`
}

// generate an enum with display values and descriptions.
type SushiEnum struct {
	Maki    string `enum:"MAKI,Rice and filling wrapped in seaweed"`
	Temaki  string `enum:"TEMAKI,Hand rolled into a cone shape"`
	Sashimi string `enum:"SASHIMI,Fish or shellfish served alone without rice"`
}
