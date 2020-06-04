package parse_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/parse"
)

func TestFindFlowFuncs(t *testing.T) {
	const expectedFuncNames = `github.com/flowdev/ea-flow-doc/parse/testdata/find_funcs: DoAccountingMagic, DoSpecialAccountingMagic, SimpleFunc, funcWithEllipsis`

	pkgs, err := parse.Dir(mustAbs(filepath.Join("testdata", "find_funcs")), false)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	pkgFuncs := parse.FindFlowFuncs(pkgs)

	actualFuncNames := funcNames(pkgFuncs)
	if actualFuncNames != expectedFuncNames {
		t.Errorf("expected functions %q but got: %q", expectedFuncNames, actualFuncNames)
	}
}

func funcNames(pkgFuncs []parse.PackageFuncs) string {
	files := make([]string, 0, len(pkgFuncs)*16)
	for _, pkgFunc := range pkgFuncs {
		names := make([]string, len(pkgFunc.Funcs))
		for i, fun := range pkgFunc.Funcs {
			names[i] = fun.Name.Name
		}
		sort.Strings(names)
		files = append(files, pkgFunc.PkgPath+": "+strings.Join(names, ", "))
	}
	sort.Strings(files)
	return strings.Join(files, " | ")
}
