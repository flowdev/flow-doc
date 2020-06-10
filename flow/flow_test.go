package flow_test

import (
	"path/filepath"
	"testing"

	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/flow"
	"github.com/flowdev/ea-flow-doc/parse"
)

func TestParse(t *testing.T) {
	root := mustAbs(filepath.Join("testdata", "functyps"))
	pkgs, err := parse.Dir(root, false)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	pkgFuncs := find.FlowFuncs(pkgs)

	flow.Parse(pkgFuncs)
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
