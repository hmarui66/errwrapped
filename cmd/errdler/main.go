package main

import (
	"github.com/hmarui66/errdler"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(errdler.Analyzer) }