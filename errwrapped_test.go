package errwrapped_test

import (
	"testing"

	"github.com/hmarui66/errwrapped"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, errwrapped.Analyzer, "a")
}
