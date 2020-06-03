package parse_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/parse"
	"golang.org/x/tools/go/packages"
)

func TestDirTree(t *testing.T) {
	specs := []struct {
		name          string
		expectedError bool
		expectedPkgs  string
	}{
		{
			name:          "happy-path",
			expectedError: false,
			expectedPkgs: "alltst: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/alltst | " +
				"alltst: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/alltst [T] | " +
				"alltst_test: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/alltst_test [T] | " +
				"apitst: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/apitst | " +
				"apitst_test: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/apitst_test [T] | " +
				"main: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path | " +
				"main: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/alltst.test [T] | " +
				"main: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/apitst.test [T] | " +
				"main: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/unittst.test [T] | " +
				"unittst: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/unittst | " +
				"unittst: github.com/flowdev/ea-flow-doc/parse/testdata/happy-path/unittst [T]",
		}, {
			name:          "error-path",
			expectedError: true,
			expectedPkgs:  "",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualPkgs, err := parse.Dir(mustAbs(filepath.Join("testdata", spec.name)), true)
			//t.Logf("err: %v, actualPkgs: %#v", err, actualPkgs)
			if spec.expectedError {
				if err != nil {
					t.Logf("received expected error: %v", err)
				} else {
					t.Error("expected to receive error but didn't get one")
				}
			} else if err != nil {
				t.Fatalf("received UNexpected error: %v", err)
			}
			actualPkgsString := packagesAsString(actualPkgs)
			if actualPkgsString != spec.expectedPkgs {
				t.Errorf("expected parsed packages %q, actual %q (len=%d)",
					spec.expectedPkgs, actualPkgsString, len(actualPkgs))
			}
		})
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
	return strings.Join(strPkgs, " | ")
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}

func isTestPkg(pkg *packages.Package) bool {
	return strings.HasSuffix(pkg.PkgPath, "_test") ||
		strings.HasSuffix(pkg.PkgPath, ".test") ||
		strings.HasSuffix(pkg.ID, ".test]") ||
		strings.HasSuffix(pkg.ID, ".test")
}
