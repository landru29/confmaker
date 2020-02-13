package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {

	var out io.Writer
	var err error

	out = os.Stdout

	inFilename := flag.String("in", "", "input filename")
	outFilename := flag.String("out", "", "output filename")

	flag.Parse()

	if *outFilename != "" {
		out, err = os.Create(*outFilename)
		if err != nil {
			panic(err)
		}
	}

	if *inFilename == "" {
		panic(fmt.Errorf("Missing go input file"))
	}

	data, err := ioutil.ReadFile(*inFilename)
	if err != nil {
		panic(err)
	}

	err = process(string(data), out)
	if err != nil {
		panic(err)
	}

}
