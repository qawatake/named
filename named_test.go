package named_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/qawatake/named"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestAnalyzer_notfound(t *testing.T) {
	t.Parallel()
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	treporter := NewAnalysisErrorReporter(1)
	analysistest.Run(treporter, testdata, named.NewAnalyzer(
		named.Deferred{
			PkgPath:  "a",
			FuncName: ".wrapper.Wrap",
			ArgPos:   0,
		},
	), "a")
	errs := treporter.Errors()
	want := named.ErrInvalidFuncName{
		FuncName: ".wrapper.Wrap",
	}
	if len(errs) != 1 {
		t.Fatalf("err expected but not found: %v", want)
	}
	if !errors.Is(errs[0], want) {
		t.Errorf("got %v, want %v", errs[0], want)
	}
}

var _ analysistest.Testing = (*analysisErrorReporter)(nil)

type analysisErrorReporter struct {
	sync.RWMutex
	errs []error
}

func NewAnalysisErrorReporter(expected int) *analysisErrorReporter {
	return &analysisErrorReporter{
		errs: make([]error, 0, expected),
	}
}

func (r *analysisErrorReporter) Errorf(format string, args ...any) {
	errs := make([]error, 0, len(args))
	for _, arg := range args {
		if err, ok := arg.(error); ok {
			errs = append(errs, err)
		}
	}
	errs = append(errs, fmt.Errorf(format, args...))
	r.Lock()
	defer r.Unlock()
	r.errs = append(r.errs, errors.Join(errs...))
}

func (r *analysisErrorReporter) Errors() []error {
	r.RLock()
	defer r.RUnlock()
	return r.errs[:]
}
