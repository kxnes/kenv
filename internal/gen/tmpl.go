// Package gen contains instruments for generation.
package gen

const (
	Get = `
// Get{{ .Title }} decorates predefined by kenv convert function to follow one code style.
func ({{ .Recv.Name }} *{{ .Recv.Type }}) Get{{ .Title }}(s string) ({{ .Field.Type }}, error) {
	return {{ .Field.Func }}(s)
}
`
	New = `
// New returns new {{ .Recv.Type }}. 
func New{{ .Recv.Type }}() *{{ .Recv.Type }} {
	{{ .Recv.Name }} := new({{ .Recv.Type }})
	{{- range $key, $value := .Data}}
	{{ $value.Recv.Name}}.{{ $value.Field.Name }} = {{ $value.Recv.Name }}.{{ $value.Field.Action }}{{ $value.Title }}("{{ $value.Field.EnvVar }}")
	{{- end}}
	return {{ .Recv.Name }}
}
`
	Want = `
// Want{{ .Title }} returns the {{ .Field.Type }} environment variable
// or default value of {{ .Field.Type }} otherwise. No checks here.
func ({{ .Recv.Name }} *{{ .Recv.Type }}) Want{{ .Title }}(key string) {{ .Field.Type }} {
	env, _ := os.LookupEnv(key)
	val, _ := {{ .Recv.Name }}.{{ .Field.Func }}(env)
	return val
}
`
	Must = `
// Must{{ .Title }} returns the {{ .Field.Type }} environment variable
// if it is exist and valid or panics otherwise.
func ({{ .Recv.Name }} *{{ .Recv.Type }}) Must{{ .Title }}(key string) {{ .Field.Type }} {
	env, ok := os.LookupEnv(key)
	if !ok {
		panic("error missing " + key)
	}

	val, err := {{ .Recv.Name }}.{{ .Field.Func }}(env)
	if err != nil {
		panic("error convert " + key)
	}

	return val
}
`
	Secret = `
// Secret{{ .Title }} returns the {{ .Field.Type }} environment variable
// if it is exist and valid or panics otherwise.
// Also Secret{{ .Title }} immediately deletes it from environment.
func ({{ .Recv.Name }} *{{ .Recv.Type }})Secret{{ .Title }}(key string) {{ .Field.Type }} {
	env, ok := os.LookupEnv(key)
	if !ok {
		panic("error missing " + key)
	}
	_ = os.Unsetenv(key)

	val, err := {{ .Recv.Name }}.{{ .Field.Func }}(env)
	if err != nil {
		panic("error convert " + key)
	}

	return val
}
`
)
