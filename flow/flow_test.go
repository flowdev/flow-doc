package flow

import (
	"path/filepath"
	"testing"

	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/parse"
)

func TestParse(t *testing.T) {
	root := mustAbs(filepath.Join("testdata", "functyps"))
	pkgs, err := parse.Dir(root, false)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	pkgFuncs := find.FlowFuncs(pkgs)

	flowDats, errs := parseAll(pkgFuncs)
	if len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}
	t.Logf("len(flowDats): %d, flowDats:", len(flowDats))
	for i, fd := range flowDats {
		t.Logf("flowDats[%d]: %s", i, fd.String())
	}
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
