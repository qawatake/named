package named

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const name = "named"
const doc = "named is ..."

func NewAnalyzer(deferred ...Deferred) *analysis.Analyzer {
	r := runner{
		deferred: deferred,
	}
	return &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		Run:  r.run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}
}

type runner struct {
	deferred []Deferred
}

type Deferred struct {
	PkgPath  string
	FuncName string
	ArgPos   int
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	m := make(map[types.Object]Deferred)
	for _, d := range r.deferred {
		obj := analysisutil.ObjectOf(pass, d.PkgPath, d.FuncName)
		if obj == nil {
			continue
		}
		m[obj] = d
	}

	inspect.WithStack(nil, func(n ast.Node, push bool, stack []ast.Node) bool {
		if stmt, ok := n.(*ast.DeferStmt); ok {
			switch f := stmt.Call.Fun.(type) {
			case *ast.Ident:
				if d, ok := m[pass.TypesInfo.ObjectOf(f)]; ok {
					for i := len(stack) - 1; i >= 0; i-- {
						n := stack[i]
						switch n := n.(type) {
						case *ast.FuncDecl:
							if !isNamedReturnValue(pass, stmt.Call.Args[d.ArgPos], n.Type.Results) {
								pass.Reportf(stmt.Call.Fun.Pos(), "%s should be called with a named return value as the %dth argument", d.FuncName, d.ArgPos+1)
							}
							return false
						case *ast.FuncLit:
							if !isNamedReturnValue(pass, stmt.Call.Args[d.ArgPos], n.Type.Results) {
								pass.Reportf(stmt.Call.Fun.Pos(), "%s should be called with a named return value as the %dth argument", d.FuncName, d.ArgPos+1)
							}
							return false
						}
					}
				}
			case *ast.SelectorExpr:
				if d, ok := m[pass.TypesInfo.ObjectOf(f.Sel)]; ok {
					for i := len(stack) - 1; i >= 0; i-- {
						n := stack[i]
						switch n := n.(type) {
						case *ast.FuncDecl:
							if !isNamedReturnValue(pass, stmt.Call.Args[d.ArgPos], n.Type.Results) {
								pass.Reportf(stmt.Call.Fun.Pos(), "%s should be called with a named return value as the %dth argument", d.FuncName, d.ArgPos+1)
							}
							return false
						case *ast.FuncLit:
							if !isNamedReturnValue(pass, stmt.Call.Args[d.ArgPos], n.Type.Results) {
								pass.Reportf(stmt.Call.Fun.Pos(), "%s should be called with a named return value as the %dth argument", d.FuncName, d.ArgPos+1)
							}
							return false
						}
					}
				}
			}
		}
		return true
	})

	return nil, nil
}

func isNamedReturnValue(pass *analysis.Pass, arg ast.Expr, fields *ast.FieldList) bool {
	unary, ok := arg.(*ast.UnaryExpr)
	if !ok {
		return false
	}
	if unary.Op != token.AND {
		return false
	}
	v, ok := unary.X.(*ast.Ident)
	if !ok {
		return false
	}
	if fields == nil {
		return false
	}
	val := pass.TypesInfo.ObjectOf(v)
	for _, field := range fields.List {
		for _, name := range field.Names {
			if val == pass.TypesInfo.ObjectOf(name) {
				return true
			}
		}
	}
	return false
}