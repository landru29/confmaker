package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type structure struct {
	name string
	decl []declaration
}

const typing = "type {{.FuncSignature}} func(*{{.StructName}}) error"

const newInstance = `
func New{{.StructName}}(config ...{{if .WithTyping}}{{.FuncSignature}}{{else}}func(*{{.StructName}}) error{{end}}) (*{{.StructName}}, error) {
	ret := &{{.StructName}}{}
	for _, c := range config {
		err := c(ret)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}
`

func (s *structure) generateStruct(withTyping bool) (string, error) {
	str := ""
	for _, d := range s.decl {

		s, err := d.generate(s.name, withTyping)
		if err != nil {
			return str, err
		}
		str = fmt.Sprintf("%s\n%s", str, s)
	}
	return str, nil
}

func (s *structure) generateTyping() (string, error) {
	type Injector struct {
		StructName    string
		FuncSignature string
	}
	inj := Injector{
		StructName:    s.name,
		FuncSignature: fmt.Sprintf("%sConfigurator", s.name),
	}
	tmpl, err := template.New("type").Parse(typing)
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

func (s *structure) generateNew(withTyping bool) (string, error) {
	type Injector struct {
		StructName    string
		FuncSignature string
		WithTyping    bool
	}

	inj := Injector{
		StructName:    s.name,
		FuncSignature: fmt.Sprintf("%sConfigurator", s.name),
		WithTyping:    withTyping,
	}
	tmpl, err := template.New("type").Parse(newInstance)
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
