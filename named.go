package named

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
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
	m := make(map[types.Object]Deferred)
	for _, d := range r.deferred {
		obj := analysisutil.ObjectOf(pass, d.PkgPath, d.FuncName)
		if obj == nil {
			continue
		}
		m[obj] = d
	}
	// unused case
	if len(m) == 0 {
		return nil, nil
	}
	fmt.Println(pass.Pkg, m)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			decl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if decl.Type.Results == nil {
				continue
			}
			returns := decl.Type.Results.List
			for _, stmt := range decl.Body.List {
				if stmt, ok := stmt.(*ast.DeferStmt); ok {
					switch f := stmt.Call.Fun.(type) {
					case *ast.Ident:
						if d, ok := m[pass.TypesInfo.ObjectOf(f)]; ok {
							if !isNamedReturnValue(pass, stmt.Call.Args[d.ArgPos], returns) {
								pass.Reportf(stmt.Call.Fun.Pos(), "%s should be called with a named return value as the %dth argument", d.FuncName, d.ArgPos+1)
							}
						}
					case *ast.SelectorExpr:
						if d, ok := m[pass.TypesInfo.ObjectOf(f.Sel)]; ok {
							if !isNamedReturnValue(pass, stmt.Call.Args[d.ArgPos], returns) {
								pass.Reportf(stmt.Call.Fun.Pos(), "%s should be called with a named return value as the %dth argument", d.FuncName, d.ArgPos+1)
							}
						}
					}
				}
			}
		}
	}
	return nil, nil
}

func isNamedReturnValue(pass *analysis.Pass, arg ast.Expr, fields []*ast.Field) bool {
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
	val := pass.TypesInfo.ObjectOf(v)
	for _, field := range fields {
		for _, name := range field.Names {
			if val == pass.TypesInfo.ObjectOf(name) {
				return true
			}
		}
	}
	return false
}
