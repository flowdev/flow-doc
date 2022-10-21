package dirs_test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/flowdev/ea-flow-doc/x/dirs"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestFindRoot(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/find-root",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"expectFindRoot": func(ts *testscript.TestScript, _ bool, args []string) {
				workDir := ts.Getenv("WORK")
				ts.Logf("$WORK: %q", workDir)
				curDir, err := os.Getwd()
				if err != nil {
					ts.Fatalf("unable to read the current directory: %v", err)
				}
				if len(args) != 3 {
					ts.Fatalf("expected 3 arguments ("+
						"givenRoot, givenIgnoreVendor, expectedRoot"+
						") but got: : %q", args)
				}
				givenIgnoreVendor, err := strconv.ParseBool(args[1])
				if err != nil {
					ts.Fatalf("the 2. argument (givenIgnoreVendor) "+
						"should be 'true' or 'false' but is: %q", args[1])
				}
				givenRoot, expectedRoot := args[0], args[2]
				expectedRoot = strings.ReplaceAll(expectedRoot, "$WORK", workDir)

				err = os.Chdir(ts.MkAbs("."))
				if err != nil {
					ts.Fatalf("unable to change the current directory: %v", err)
				}

				actualRoot, err := dirs.FindRoot(givenRoot, givenIgnoreVendor)
				os.Chdir(curDir)
				if err != nil {
					ts.Fatalf("expected no error but got: %v", err)
				}
				ts.Logf("expectedRoot: %q, actualRoot: %q", expectedRoot, actualRoot)
				if actualRoot != expectedRoot {
					ts.Fatalf("expected project root %q, got: %q",
						expectedRoot, actualRoot)
				}
			},
		},
		TestWork: false,
	})
}
