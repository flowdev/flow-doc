package dirs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flowdev/ea-flow-doc/x/dirs"
)

func TestFindRoot(t *testing.T) {
	testDataDir := mustAbs(filepath.Join("testdata", "find-root"))
	givenStartDir := filepath.Join("in", "some", "subdir")
	specs := []struct {
		name              string
		givenRoot         string
		givenIgnoreVendor bool
		expectedRoot      string
	}{
		{
			name:              "go-mod",
			givenRoot:         "",
			givenIgnoreVendor: false,
			expectedRoot:      filepath.Join(testDataDir, "go-mod"),
		}, {
			name:              "given-root",
			givenRoot:         "/my/given/root/dir",
			givenIgnoreVendor: false,
			expectedRoot:      "/my/given/root/dir",
		}, {
			name:              "vendor-dir",
			givenRoot:         "",
			givenIgnoreVendor: false,
			expectedRoot:      filepath.Join(testDataDir, "vendor-dir"),
		}, {
			name:              "ignore-vendor",
			givenRoot:         "",
			givenIgnoreVendor: true,
			expectedRoot:      filepath.Join(testDataDir, "ignore-vendor"),
		},
	}

	initDir := mustAbs(".")
	t.Cleanup(func() {
		mustChdir(initDir)
	})
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			mustChdir(filepath.Join(testDataDir, spec.name, givenStartDir))

			actualRoot, err := dirs.FindRoot(spec.givenRoot, spec.givenIgnoreVendor)
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
			if actualRoot != spec.expectedRoot {
				t.Errorf("expected project root %q, actual %q",
					spec.expectedRoot, actualRoot)
			}
		})
	}
}

func mustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		panic(err.Error())
	}
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
