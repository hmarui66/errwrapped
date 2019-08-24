package errdler

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	wrapper string
)

func init() {
	Analyzer.Flags.StringVar(&wrapper, "wrapper", `errors`, "name of error wrapper")
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
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fd, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}
		errIdx, exists := getErrorIdx(fd)
		if !exists {
			return
		}

		if fd.Body == nil || len(fd.Body.List) == 0 {
			return
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

			sel, ok := cal.Fun.(*ast.SelectorExpr)
			if !ok {
				detected = append(detected, ret)
				return false
			}

			if id, ok := sel.X.(*ast.Ident); !ok || id.Name != wrapper {
				detected = append(detected, ret)
				return false
			}

			return true
		})

		for _, n := range detected {
			//ast.Print(pass.Fset, n)
			pass.Reportf(n.Pos(), "unhandled error found")
		}
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
