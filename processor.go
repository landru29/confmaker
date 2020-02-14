package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
)

func process(src string, out io.Writer, withTyping bool) error {
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

			if withTyping {

				str, err := st.generateTyping()
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "\n%s", str)
			}

			str, err := st.generateNew(withTyping)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "\n%s", str)

			str, err = st.generateStruct(withTyping)
			if err != nil {
				return err
			}

			fmt.Fprintf(out, str)
		}
	}

	return err
}
