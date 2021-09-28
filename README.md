# go-ultra-enum


`go-ultra-enum` is an enum generator for Go. It is inspired by the powerful enum types found in Java. `go-ultra-enum` has the following capabilities

* Reference and compare enums using values
* Support multi value type, such as int, bool.

## Install

From a github release


```bash
go get -u github.com/drinks5/go-ultra-enum
```

## Example

To define an enum, create a `struct` with the suffix `Enum`. You can define a display value in the `struct` tag. Adding a hyphen will assign the field name to the display value.

You can then generate the enum as follows.

```go
//go:generate go-ultra-enum -file=$GOFILE

// generate an enum with display values. The display values are used for JSON serialization/deserialization
type ColorEnum struct {
	Red       int `enum:"2"`
	LightBlue int `enum:"1"`
}

type GeoEnum struct {
	Lat int64 `enum:"2"`
	Lon int64 `enum:"1"`
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

```

When a description is defined the json is serialized as follows (not yet implemented)

```json
{
  "sushi": {
    "name": "SASHIMI",
    "description": "Fish or shellfish served alone without rice"
  }
}
```

## Consumer api

The generated code would yield the following consumer api

### Create an enum value

```go
a := Red // OR var a Color = Red
```

### Create an enum from a factory method

```go
var name Color = NewColor("RED")
```

### Get the display value

```go
var name string = a.Name() // "RED"
```

### Get all display values

```go
var names []string = ColorNames() // []string{"RED", "BLUE"}
```

### Get all values

```go
var values []Color = ColorValues() // []string{Red, Blue}
```

### Pass as an error

Enums implement `Error() string` which means they can be passed as errors.

```go
var a error = Red
```

## Developing

```bash
go build main.go
go generate
go test .
```

OR

```bash
make test
```
