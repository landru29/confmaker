package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type declaration struct {
	name     string
	typeName string
}

const configurator = `
func With{{.Field}}({{.AttrName}} {{.FieldType}}) {{if .WithTyping}}{{.FuncSignature}}{{else}}func(*{{.StructName}}) error{{end}} {
	return func({{.Variable}} *{{.StructName}}) error {
		{{.Variable}}.{{.FieldName}} = {{.AttrName}}
		return nil
	}
}
`

func (d declaration) generate(name string, withTyping bool) (string, error) {
	type Injector struct {
		Variable      string
		StructName    string
		Field         string
		FieldName     string
		AttrName      string
		FieldType     string
		FuncSignature string
		WithTyping    bool
	}

	attrName := strings.ToLower(d.name[:1])
	variable := strings.ToLower(name[:1])

	if variable == attrName {
		if len(d.name) > 1 {
			attrName = strings.ToLower(d.name[:2])
		} else {
			attrName = fmt.Sprintf("%sa", attrName)
		}
	}

	inj := Injector{
		Variable:      variable,
		StructName:    name,
		Field:         strings.Title(d.name),
		FieldName:     d.name,
		AttrName:      attrName,
		FieldType:     d.typeName,
		FuncSignature: fmt.Sprintf("%sConfigurator", name),
		WithTyping:    withTyping,
	}
	tmpl, err := template.New("decl").Parse(configurator)
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
