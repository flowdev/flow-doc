package find_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/parse"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestFlowFuncsAndTests(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"findFlowFuncs": findFlowFuncs,
			"findFlowTests": findFlowTests,
		},
		// TestWork: true,
	})
}

func findFlowFuncs(ts *testscript.TestScript, _ bool, args []string) {
	workDir := ts.Getenv("WORK")
	funcFile := workDir + "/flowFuncs.actual"

	pkgs, err := parse.Dir(workDir, true)
	if err != nil {
		ts.Fatalf("received unexpected error: %v", err)
	}

	flowFuncs := find.FlowFuncs(pkgs)

	actualFuncNames := funcNames(flowFuncs, workDir)
	err = os.WriteFile(funcFile, []byte(actualFuncNames+"\n\n"), 0666)
	if err != nil {
		ts.Fatalf("ERROR: Unable to write file '%s': %v", funcFile, err)
	}
}

func findFlowTests(ts *testscript.TestScript, _ bool, args []string) {
	workDir := ts.Getenv("WORK")
	testsFile := workDir + "/flowTests.actual"

	pkgs, err := parse.Dir(workDir, true)
	if err != nil {
		ts.Fatalf("received unexpected error: %v", err)
	}

	flowTests := find.FlowTests(pkgs)

	actualTestNames := funcNames(flowTests, workDir)
	err = os.WriteFile(testsFile, []byte(actualTestNames+"\n\n"), 0666)
	if err != nil {
		ts.Fatalf("ERROR: Unable to write file '%s': %v", testsFile, err)
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
	return strings.Join(names, "\n")
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
