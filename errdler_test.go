package errdler_test

import (
	"testing"

	"github.com/hmarui66/errdler"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, errdler.Analyzer, "a")
}