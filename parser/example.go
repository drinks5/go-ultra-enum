package parser

//go:generate go-ultra-enum
type ColorEnum struct {
	Red   int `enum:"1"`
	Blue  int `enum:"2"`
	Grenn int `enum:"3"`
}

type CountryEnum struct {
	China    int64 `enum:"1"`
	America  int64 `enum:"2"`
	Sinapore int64 `enum:"3"`
}

// generate an enum with default display values. The display values are set to the field names, e.g. `On` and `Off`
type StatusEnum struct {
	On  bool `enum:"true"`
	Off bool `enum:"false"`
}

// generate an enum with display values and descriptions.
type SushiEnum struct {
	Maki    string `enum:"MAKI,Rice and filling wrapped in seaweed"`
	Temaki  string `enum:"TEMAKI,Hand rolled into a cone shape"`
	Sashimi string `enum:"SASHIMI,Fish or shellfish served alone without rice"`
}
