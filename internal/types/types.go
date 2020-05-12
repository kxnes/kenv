// Package types contains types for between parser and code generator.
package types

const (
	Want   = "-" // It will be nice if environment variable exist.
	Must   = "!" // Environment variable must be exist and fine.
	Secret = "*" // Environment variable must be exist and fine
	//               and will be deleted from environment after init.
)

type (
	// Tag describes the env tag values. Tag format: `env:"[ENV]ACTION"`.
	//  - ENV    - explicit environment variable name or `Field.Name`.ToUpper().
	//  - ACTION - how to get environment variable (see above).
	Tag struct {
		EnvVar string // Name of environment variable.
		Action string // Action on environment variable [* ! -].
	}

	// Field describes the struct field.
	Field struct {
		Name     string // Name of field.
		Type     string // Name of field type.
		ConvFunc string // Convert function for conversion from `string` to `Type`.
		EnvVar   string // Name of environment variable.
		Action   string // Action on environment variable [* ! -].
	}
)
