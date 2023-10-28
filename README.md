# named

Linter `named` ensures a target function to be called with a named return value.

A typical use case is to prevent the misuse of a error wrapping function such as [derrors.Wrap](https://github.com/golang/pkgsite/blob/5f0513d53cff8382238b5f8c78e8317d2b4ad06d/internal/derrors/derrors.go#L240), which does not allow the resulted error to be unwrapped.
Function `Bad` below fails to wrap the error simply because it doesn't use the named return value as an argument.

```go
func Good() (err error) {
  err = DoSomething()
  defer Wrap(&err, "wrapped!!") // ok because a named return value is passed.
  return fmt.Errorf("error from Good: %w", err)
}

func Bad() error {
  err := DoSomething()
  defer Wrap(&err, "wrapped!!") // <- err is not a named return value.
  return fmt.Errorf("error from Bad: %w", err)
}

func Wrap(errp *error, msg string) {
  if *errp == nil {
    return
  }
  *errp = fmt.Errorf("%s: %w", msg, *errp)
}

func DoSomething() error {
  return errors.New("original error")
}

func main() {
  err := Good()
  fmt.Println(err)
  err = Bad()
  fmt.Println(err)
  // Output:
  // wrapped!!: error from Good: original error
  // error from Bad: original error
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
