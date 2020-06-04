package parse_test

import (
	"go/ast"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/parse"
)

func TestFindFlowFuncs(t *testing.T) {
	const expectedFuncNames = `DoAccountingMagic | DoSpecialAccountingMagic | SimpleFunc | funcWithEllipsis`

	pkgs, err := parse.Dir(mustAbs(filepath.Join("testdata", "find_flows")), false)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	flowFuncs := parse.FindFlowFuncs(pkgs)

	actualFuncNames := funcNames(flowFuncs)
	if actualFuncNames != expectedFuncNames {
		t.Errorf("expected functions %q but got: %q", expectedFuncNames, actualFuncNames)
	}
}

func funcNames(funcs []*ast.FuncDecl) string {
	names := make([]string, len(funcs))
	for i, fun := range funcs {
		names[i] = fun.Name.Name
	}
	sort.Strings(names)
	return strings.Join(names, " | ")
}
