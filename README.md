kenv
====

> Experiments with GoLang code generation.

Example
-------

```bash
$ go get github.com/kxnes/kenv
```

```go
package main

import "fmt"

type Addr struct {
	Host string `env:"!"`
	Port string `env:"!"`
}

type Database struct {
	Addr     Addr
	Name     string `env:"DB_NAME-"`
	User     string `env:"DB_USER!"`
	Password int    `env:"DB_PASSWORD*"`
}

//go:generate kenv -t Environment -f ./env.go
type Environment struct {
	Database Database
	UID      int `env:"UID!"`
}

func main() {
	fmt.Println(NewEnvironment()) // will be generated
}
```
