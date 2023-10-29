package named_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/qawatake/named"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, named.NewAnalyzer(
		named.Deferred{
			PkgPath:  "a",
			FuncName: "Wrap",
			ArgPos:   0,
		},
		named.Deferred{
			PkgPath:  "a/b",
			FuncName: "Wrap",
			ArgPos:   0,
		},
		named.Deferred{
			PkgPath:  "a",
			FuncName: "WrapAny",
			ArgPos:   0,
		},
		named.Deferred{
			PkgPath:  "a",
			FuncName: "wrapper.Wrap",
			ArgPos:   0,
		},
	), "a/...")

}

func TestAnalyzer_pkgname_is_different_from_pkgpath(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, named.NewAnalyzer(
		named.Deferred{
			PkgPath:  "github.com/qawatake/a",
			FuncName: "Wrap",
			ArgPos:   0,
		},
		named.Deferred{
			PkgPath:  "github.com/qawatake/a",
			FuncName: "wrapper.Wrap",
			ArgPos:   0,
		},
	), "github.com/qawatake/a/...")
}
