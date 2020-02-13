package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"text/template"
)

const configurator = `
func With{{.FieldName}}({{.AttrName}} {{.FieldType}}) func(*{{.StructName}}) error {
	return func({{.Variable}} *{{.StructName}}) error {
		{{.Variable}}.{{.FieldName}} = {{.AttrName}}
		return nil
	}
}
`

type goSource string

type declaration struct {
	name     string
	typeName string
}

type structure struct {
	name string
	decl []declaration
}

func (s *structure) generateStruct() (string, error) {
	str := ""
	for _, d := range s.decl {
		s, err := d.generate(s.name)
		if err != nil {
			return str, err
		}
		str = fmt.Sprintf("%s\n%s", str, s)
	}
	return str, nil
}

func (d declaration) generate(name string) (string, error) {

	type Injector struct {
		Variable   string
		StructName string
		FieldName  string
		AttrName   string
		FieldType  string
	}
	inj := Injector{
		Variable:   strings.ToLower(name[:1]),
		StructName: name,
		FieldName:  d.name,
		AttrName:   strings.ToLower(d.name[:1]),
		FieldType:  d.typeName,
	}
	tmpl, err := template.New("test").Parse(configurator)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")

	err = tmpl.Execute(buf, inj)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func process(src string, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(src), 0)

	if err != nil {
		panic(err)
	}

	start := f.Name.Pos() - 1
	end := f.Name.End() - 1
	pkName := src[start:end]

	fmt.Fprintf(out, "package %s\n", pkName)

	for _, d := range f.Decls {
		st := goSource(src).structures(d)

		if st != nil {
			str, err := st.generateStruct()
			if err != nil {
				return err
			}

			fmt.Fprintf(out, str)
		}
	}

	return err
}

func (src goSource) structures(d ast.Decl) *structure {
	var res *structure
	typeDecl := d.(*ast.GenDecl)

	if len(typeDecl.Specs) == 1 {

		if spec, ok := typeDecl.Specs[0].(*ast.TypeSpec); ok {

			if structDecl, ok := spec.Type.(*ast.StructType); ok {
				res = &structure{}

				start := spec.Name.Pos() - 1
				end := spec.Name.End() - 1
				res.name = string(src)[start:end]

				res.decl = src.fields(structDecl.Fields.List)
			}
		}
	}
	return res
}

func (src goSource) fields(fields []*ast.Field) []declaration {
	res := []declaration{}

	for _, field := range fields {
		typeName := src.fieldType(field)
		for _, name := range src.fieldName(field) {
			res = append(res, declaration{name: name, typeName: typeName})
		}
	}
	return res
}

func (src goSource) fieldType(field *ast.Field) string {
	typeExpr := field.Type

	start := typeExpr.Pos() - 1
	end := typeExpr.End() - 1

	// grab it in source
	return string(src)[start:end]
}

func (src goSource) fieldName(field *ast.Field) []string {
	res := []string{}

	for _, typeExpr := range field.Names {

		start := typeExpr.Pos() - 1
		end := typeExpr.End() - 1

		// grab it in source
		res = append(res, string(src)[start:end])
	}
	return res
}
