// Package conv contains internal checks of public conv package.
package conv

import (
	"fmt"

	kenv "github.com/kxnes/kenv/conv"
)

// GetPredefinedImport uses if need import of predefined convert functions.
// Holds a copy of import above. MUST be set manual.
func GetPredefinedImport() []string {
	return []string{
		`"os"`,
		`kenv "github.com/kxnes/kenv/conv"`,
	}
}

// GetPredefinedConv returns all predefined convert functions.
// Holds a copy of functions below. MUST be set manual.
func GetPredefinedConv() map[string]string {
	return map[string]string{
		"bool": "kenv.Bool",

		"uint":   "kenv.Uint",
		"uint8":  "kenv.Uint8",
		"uint16": "kenv.Uint16",
		"uint32": "kenv.Uint32",
		"uint64": "kenv.Uint64",

		"int":   "kenv.Int",
		"int8":  "kenv.Int8",
		"int16": "kenv.Int16",
		"int32": "kenv.Int32",
		"int64": "kenv.Int64",

		"string": "kenv.String",
		"rune":   "kenv.Rune",

		"float32": "kenv.Float32",
		"float64": "kenv.Float64",
	}
}

// init uses to ensure package "pkg/conv" integrity and interfaces.
func init() {
	must(kenv.Bool("true"))

	must(kenv.Uint("1234"))
	must(kenv.Uint8("123"))
	must(kenv.Uint16("12"))
	must(kenv.Uint32("12"))
	must(kenv.Uint64("12"))

	must(kenv.Int("1234"))
	must(kenv.Int8("123"))
	must(kenv.Int16("12"))
	must(kenv.Int32("12"))
	must(kenv.Int64("12"))

	must(kenv.String("s"))
	must(kenv.Rune("Я"))

	must(kenv.Float32("1.2"))
	must(kenv.Float64("1.2"))
}

// must decorates `init` calls.
func must(_ interface{}, err error) {
	if err != nil {
		panic(fmt.Errorf("integrity error: %w", err))
	}
}
