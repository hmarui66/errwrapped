package errwrapped

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type options []string

func (o *options) String() string {
	return fmt.Sprint(*o)
}

func (o *options) Set(value string) error {
	if len(*o) > 0 {
		return errors.New("option flag already set")
	}
	for _, opt := range strings.Split(value, ",") {
		*o = append(*o, opt)
	}
	return nil
}

func (o *options) Contains(value string) bool {
	if o == nil {
		return false
	}
	for _, opt := range *o {
		if strings.Contains(value, opt) {
			return true
		}
	}

	return false
}

var (
	wrapperFlag   options
	ignoreFlag    options
	ignoreOneline bool
)

func init() {
	Analyzer.Flags.Var(&wrapperFlag, "wrapper", "comma-separated list of error wrapper name")
	Analyzer.Flags.Var(&ignoreFlag, "ignore", "comma-separated list of ignoring file name suffix")
	Analyzer.Flags.BoolVar(&ignoreOneline, "ignore-oneline", false, "ignore one line function")
}

var Analyzer = &analysis.Analyzer{
	Name: "errdler",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "errdler is ..."

func run(pass *analysis.Pass) (interface{}, error) {
	if len(wrapperFlag) == 0 {
		if err := wrapperFlag.Set("errors"); err != nil {
			return nil, err
		}
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Nodes(nodeFilter, func(n ast.Node, push bool) bool {
		filename := pass.Fset.File(n.Pos()).Name()
		if ignoreFlag.Contains(filename) {
			return false
		}

		fd, ok := n.(*ast.FuncDecl)
		if !ok {
			return false
		}
		errIdx, exists := getErrorIdx(fd)
		if !exists {
			return false
		}

		if fd.Body == nil || len(fd.Body.List) == 0 {
			return false
		}

		if ignoreOneline && len(fd.Body.List) == 1 {
			return false
		}

		var detected []*ast.ReturnStmt
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			ret, ok := n.(*ast.ReturnStmt)
			if !ok {
				return true
			}

			if len(ret.Results)-1 < errIdx {
				detected = append(detected, ret)
				return false
			}

			errRes := ret.Results[errIdx]

			if id, ok := errRes.(*ast.Ident); ok && id.Name == "nil" {
				return false
			}

			cal, ok := errRes.(*ast.CallExpr)
			if !ok {
				detected = append(detected, ret)
				return false
			}

			if id, ok := cal.Fun.(*ast.Ident); ok && wrapperFlag.Contains(id.Name) {
				return true
			}

			sel, ok := cal.Fun.(*ast.SelectorExpr)
			if !ok {
				detected = append(detected, ret)
				return false
			}

			if id, ok := sel.X.(*ast.Ident); !ok || !wrapperFlag.Contains(id.Name) {
				detected = append(detected, ret)
				return false
			}

			return true
		})

		for _, n := range detected {
			//ast.Print(pass.Fset, n)
			pass.Reportf(n.Pos(), "unwrapped error found")
		}

		return false
	})

	return nil, nil
}

func getErrorIdx(fd *ast.FuncDecl) (int, bool) {
	if fd.Type == nil ||
		fd.Type.Results == nil ||
		len(fd.Type.Results.List) == 0 {
		return 0, false
	}

	for i, fld := range fd.Type.Results.List {
		if fld.Type == nil {
			continue
		}
		typ, ok := fld.Type.(*ast.Ident)
		if !ok || typ.Name != `error` {
			continue
		}

		return i, true
	}

	return 0, false
}
