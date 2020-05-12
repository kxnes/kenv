package main

import (
	"kenv/kenv/internal/gen"
	"kenv/kenv/internal/parser"
)

func main() {
	var (
		target   = "Environment"
		filename = "/home/kxnes/tmp/kenv/env/env.go"
	)

	p := parser.New(target, filename)
	err := gen.CodeGen(p)
	if err != nil {
		panic(err)
	}

	//p.Overview()
}
