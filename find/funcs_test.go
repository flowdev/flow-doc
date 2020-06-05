package find_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/parse"
)

func TestFlowFuncs(t *testing.T) {
	const expectedFuncNames = `funcs.go: DoAccountingMagic | funcs.go: DoSpecialAccountingMagic | funcs.go: SimpleFunc | funcs.go: funcWithEllipsis`

	root := mustAbs(filepath.Join("testdata", "funcs"))
	pkgs, err := parse.Dir(root, false)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	pkgFuncs := find.FlowFuncs(pkgs)

	actualFuncNames := funcNames(pkgFuncs, root)
	if actualFuncNames != expectedFuncNames {
		t.Errorf("expected functions %q but got: %q", expectedFuncNames, actualFuncNames)
	}
}

func TestFlowTests(t *testing.T) {
	const expectedFuncNames = `x/tool/tool_test.go: TestTool | x/tool2/tool2_test.go: TestTool2`

	root := mustAbs(filepath.Join("testdata", "funcs"))
	pkgs, err := parse.Dir(root, true)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	pkgFuncs := find.FlowTests(pkgs)

	actualFuncNames := funcNames(pkgFuncs, root)
	if actualFuncNames != expectedFuncNames {
		t.Errorf("expected functions %q but got: %q", expectedFuncNames, actualFuncNames)
	}
}

func funcNames(pkgFuncs []find.PackageFuncs, root string) string {
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

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
