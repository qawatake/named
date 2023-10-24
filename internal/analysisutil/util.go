package analysisutil

import (
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
)

// copied and modified from https://github.com/gostaticanalysis/analysisutil/blob/ccfdecf515f47e636ba164ce0e5f26810eaf8747/types.go#L18
// ObjectOf returns types.Object by given name in the package.
func ObjectOf(pass *analysis.Pass, pkg, name string) types.Object {
	obj := analysisutil.LookupFromImports(pass.Pkg.Imports(), pkg, name)
	if obj != nil {
		return obj
	}
	if analysisutil.RemoveVendor(pass.Pkg.Path()) != analysisutil.RemoveVendor(pkg) {
		return nil
	}
	return pass.Pkg.Scope().Lookup(name)
}

// copied and modified from https://github.com/gostaticanalysis/analysisutil/blob/ccfdecf515f47e636ba164ce0e5f26810eaf8747/types.go#L31
// TypeOf returns types.Type by given name in the package.
// TypeOf accepts pointer types such as *T.
func TypeOf(pass *analysis.Pass, pkg, name string) types.Type {
	if name == "" {
		return nil
	}

	if name[0] == '*' {
		obj := TypeOf(pass, pkg, name[1:])
		if obj == nil {
			return nil
		}
		return types.NewPointer(obj)
	}

	obj := ObjectOf(pass, pkg, name)
	if obj == nil {
		return nil
	}

	return obj.Type()
}

func MethodOf(typ types.Type, name string) *types.Func {
	return analysisutil.MethodOf(typ, name)
}
