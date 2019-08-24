package main

import (
	"github.com/hmarui66/errwrapped"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(errwrapped.Analyzer) }
