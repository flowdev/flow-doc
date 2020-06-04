package parse_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/parse"
)

func TestFindFlowFuncs(t *testing.T) {
	const expectedFuncNames = `funcs.go: DoAccountingMagic | funcs.go: DoSpecialAccountingMagic | funcs.go: SimpleFunc | funcs.go: funcWithEllipsis`

	root := mustAbs(filepath.Join("testdata", "find_funcs"))
	pkgs, err := parse.Dir(root, false)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	pkgFuncs := parse.FindFlowFuncs(pkgs)

	actualFuncNames := funcNames(pkgFuncs, root)
	if actualFuncNames != expectedFuncNames {
		t.Errorf("expected functions %q but got: %q", expectedFuncNames, actualFuncNames)
	}
}

func TestFindFlowTests(t *testing.T) {
	const expectedFuncNames = `x/tool/tool_test.go: TestTool | x/tool2/tool2_test.go: TestTool2`

	root := mustAbs(filepath.Join("testdata", "find_funcs"))
	pkgs, err := parse.Dir(root, true)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	pkgFuncs := parse.FindFlowTests(pkgs)

	actualFuncNames := funcNames(pkgFuncs, root)
	if actualFuncNames != expectedFuncNames {
		t.Errorf("expected functions %q but got: %q", expectedFuncNames, actualFuncNames)
	}
}

func funcNames(pkgFuncs []parse.PackageFuncs, root string) string {
	names := make([]string, 0, 4096)
	for _, pkgFunc := range pkgFuncs {
		for _, fun := range pkgFunc.Funcs {
			names = append(
				names,
				relativeFileName(pkgFunc.Fset.Position(fun.Name.NamePos).String(), root)+": "+fun.Name.Name)
		}
	}
	sort.Strings(names)
	return strings.Join(names, " | ")
}

func relativeFileName(fname, root string) string {
	fname = strings.SplitN(fname, ":", 2)[0]
	if len(fname) <= len(root) {
		return fname
	}
	fname = fname[len(root):]

	if fname[0] == '/' {
		return fname[1:]
	}
	return fname
}
