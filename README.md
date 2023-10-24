# named

Linter `named` ensures a target function to be called with a named return value.

```go
func Good() (err error) {
  defer Wrap(&err, "x") // ok because a named return value is passed.
  return nil
}

func Bad() error {
  err := fmt.Errorf("x")
  defer Wrap(&err, "x") // <- err is not a named return value.
  return err
}

func Wrap(errp *error, msg string) {
  if errp == nil || *errp == nil {
    return
  }
  *errp = fmt.Errorf("%s: %w", msg, *errp)
}
```

## How to use

Build your `named` binary by writing `main.go` like below.

```go
package main

import (
  "github.com/qawatake/named"
  "golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
  unitchecker.Main(
    named.NewAnalyzer(
      named.Deferred{
        PkgPath:  "pkg/in/which/target/func/is/defined",
        FuncName: "Wrap",
        ArgPos:   0,
      },
    ),
  )
}
```

Then, run `go vet` with your `named` binary.

```sh
go vet -vettool=/path/to/your/named ./...
```

**analysisutilのObjectOfをカスタマイズしよう**
