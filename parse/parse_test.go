package parse_test

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/parse"
	"github.com/rogpeppe/go-internal/testscript"
	"golang.org/x/tools/go/packages"
)

func TestDir(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"parseDir": parseDir,
		},
		// TestWork: true,
	})
}

func parseDir(ts *testscript.TestScript, _ bool, args []string) {
	workDir := ts.Getenv("WORK")
	pkgFile := workDir + "/packages.actual"
	output := ""

	actualPkgs, err := parse.Dir(workDir, true)
	ts.Logf("err: %v, actualPkgs: %#v", err, actualPkgs)
	if err != nil {
		ts.Logf("err: %v", err)
		output = "error: true"
	} else {
		output = "error: false\n" + packagesAsString(actualPkgs)
	}

	err = os.WriteFile(pkgFile, []byte(output+"\n\n"), 0666)
	if err != nil {
		ts.Fatalf("ERROR: Unable to write file '%s': %v", pkgFile, err)
	}
}

func packagesAsString(pkgs []*packages.Package) string {
	strPkgs := make([]string, len(pkgs))

	for i, p := range pkgs {
		strPkgs[i] = p.Name + ": " + p.PkgPath
		if isTestPkg(p) {
			strPkgs[i] += " [T]"
		}
	}
	sort.Strings(strPkgs)
	return strings.Join(strPkgs, "\n")
}

func isTestPkg(pkg *packages.Package) bool {
	return strings.HasSuffix(pkg.PkgPath, "_test") ||
		strings.HasSuffix(pkg.PkgPath, ".test") ||
		strings.HasSuffix(pkg.ID, ".test]") ||
		strings.HasSuffix(pkg.ID, ".test")
}
