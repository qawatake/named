package main

import (
	"github.com/qawatake/named"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		named.NewAnalyzer(
			named.Deferred{
				PkgPath:  "github.com/qawatake/named/internal/example",
				FuncName: "Wrap",
				ArgPos:   0,
			},
		),
	)
}
