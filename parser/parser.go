package parser

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	Name  string
	files []*File
}
type File struct {
	name string
	file *ast.File // Parsed AST.
}
type Enum struct {
	OriginalName string
	NewName      string
	Tpe          string
	Elements     []EnumElement
}

type EnumElement struct {
	Value       string
	Name        string
	Description string
	Tpe         string
}
type Parser struct {
	Pkg *Package     // Package we are scanning.
	buf bytes.Buffer // Accumulated output.
}

func (p *Parser) Packages(pattern string) {
	cfg := &packages.Config{
		Mode:  packages.LoadTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	p.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (p *Parser) addPackage(pkg *packages.Package) {
	p.Pkg = &Package{
		Name:  pkg.Name,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		p.Pkg.files[i] = &File{
			file: file,
			name: pkg.GoFiles[i],
		}
	}
}

func (p *Parser) Generate() {
	fset := token.NewFileSet()
	parserMode := parser.ParseComments
	var fileAst *ast.File
	var err error
	var enums []*Enum
	for _, file := range p.Pkg.files {
		fileAst, err = parser.ParseFile(fset, file.name, nil, parserMode)
		if err != nil {
			log.Fatal(err)
		}
		enums = append(enums, getEnumsFromFile(fileAst.Decls)...)
	}
	for _, enum := range enums {
		p.Render(Tpl, enum)
	}
}
func (p *Parser) Sink(dir string) {
	src := p.format()

	// Write to file.
	outputName := dir
	if outputName == "" {
		outputName = filepath.Join(dir, "enumer.go")
	}
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
func (p *Parser) format() []byte {
	src, err := format.Source(p.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return p.buf.Bytes()
	}
	return src
}
func getEnumsFromFile(specs []ast.Decl) (enums []*Enum) {
	for _, spec := range specs {
		decl, ok := spec.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			continue
		}
		for _, spec := range decl.Specs {
			vspec := spec.(*ast.TypeSpec)
			structType, ok := vspec.Type.(*ast.StructType)
			if !ok || structType.Fields == nil || !strings.Contains(vspec.Name.Name, "Enum") {
				continue
			}
			elements := getElementsFromFields(structType.Fields.List)
			if len(elements) != 0 {
				e := &Enum{
					OriginalName: vspec.Name.Name,
					Tpe:          elements[0].Tpe,
					NewName:      strings.Replace(vspec.Name.Name, "Enum", "", -1),
					Elements:     elements,
				}
				enums = append(enums, e)
			}
		}
	}
	return
}
func getElementsFromFields(fields []*ast.Field) []EnumElement {
	elements := make([]EnumElement, 0)
	for _, field := range fields {
		if field.Tag == nil || !strings.HasPrefix(field.Tag.Value, "`enum:") {
			continue
		}
		if len(field.Names) == 0 {
			continue
		}
		value, description := parseEnumStructTag(field.Tag.Value)
		if value == "-" {
			value = field.Names[0].Name
		}

		// grab it in source
		tpe := field.Type.(*ast.Ident).Name
		elements = append(elements, EnumElement{
			Value:       value,
			Name:        field.Names[0].Name,
			Description: description,
			Tpe:         tpe,
		})
	}
	return elements
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
func (p *Parser) Render(tmpl string, model interface{}) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"LcFirst": LcFirst,
	}
	t := template.Must(template.New(tmpl).Funcs(funcMap).Parse(tmpl))
	err := t.Execute(&p.buf, model)
	if err != nil {
		log.Fatal("Execute: ", err)
		return
	}
}

func parseEnumStructTag(content string) (string, string) {
	if value, ok := parseStructTag(content, "`enum"); ok {
		splits := strings.Split(value, ",")
		name := splits[0]
		var description string
		if len(splits) > 1 {
			description = splits[1]
		}
		return name, description
	}
	log.Fatal("enum struct tag did not contain name")
	return "", ""
}

func parseStructTag(tag string, key string) (value string, ok bool) {
	for tag != "" {
		// Skip leading space.q
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
