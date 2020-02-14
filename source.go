package main

import "go/ast"

type goSource string

func (src goSource) structures(d ast.Decl) *structure {
	var res *structure
	if typeDecl, ok := d.(*ast.GenDecl); ok {

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
