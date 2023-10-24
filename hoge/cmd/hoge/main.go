package main

import (
	"hoge"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(hoge.Analyzer) }
