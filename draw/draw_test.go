package draw_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flowdev/ea-flow-doc/draw"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestFromFlowData(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"drawBigTestFlowData": drawBigTestFlowData,
		},
		// TestWork: true,
	})
}

func drawBigTestFlowData(ts *testscript.TestScript, _ bool, args []string) {
	workDir := ts.Getenv("WORK")
	funcFile := filepath.Join(workDir, "flow.actual")

	gotBytes, err := draw.FromFlowData(BigTestFlowData)
	if err != nil {
		ts.Fatalf("unexpected error: %s", err)
	}

	err = os.WriteFile(funcFile, gotBytes, 0666)
	if err != nil {
		ts.Fatalf("unable to write file %q: %v", funcFile, err)
	}
}
