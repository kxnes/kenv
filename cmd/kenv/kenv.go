package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kxnes/kenv/internal/gen"
	"github.com/kxnes/kenv/internal/parser"
)

const errCode = 1

func missing(f, v string) {
	if v != "" {
		return
	}

	flag.Usage()
	fmt.Println("missing " + f)
	os.Exit(errCode)
}

func main() {
	var (
		target   string
		filename string
	)

	flag.StringVar(&target, "t", "", "target struct describes environment")
	flag.StringVar(&filename, "f", "", "target struct filename")
	flag.Parse()

	missing("target", target)
	missing("filename", filename)

	err := gen.CodeGen(parser.New(target, filename))
	if err != nil {
		fmt.Println(err)
		os.Exit(errCode)
	}
}
