package draw2_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/flowdev/ea-flow-doc/draw2"
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

	if len(args) != 2 {
		ts.Fatalf("expected 2 args (splitMode and darkMode), got: %q", args)
	}
	splitMode, err := strconv.ParseBool(args[0])
	if err != nil {
		ts.Fatalf("expected boolean for splitMode, got: %q; err: %v", args[0], err)
	}
	darkMode, err := strconv.ParseBool(args[1])
	if err != nil {
		ts.Fatalf("expected boolean for darkMode, got: %q; err: %v", args[1], err)
	}
	mdFile := "markdown-" + args[0] + "-" + args[1] + ".actual"

	flowMode := draw2.FlowModeNoLinks
	if splitMode {
		flowMode = draw2.FlowModeSVGLinks
	}

	BigTestFlowData.ChangeConfig("bigTestFlow", flowMode, 1500, darkMode)
	svgContents, mdContent, err := BigTestFlowData.Draw()
	if err != nil {
		ts.Fatalf("unexpected error: %s", err)
	}

	for fnam, fcontent := range svgContents {
		workFNam := filepath.Join(workDir, fnam)
		err = os.WriteFile(workFNam, fcontent, 0666)
		if err != nil {
			ts.Fatalf("unable to write file %q: %v", workFNam, err)
		}
		err = os.WriteFile(fnam, fcontent, 0666)
		if err != nil {
			ts.Fatalf("unable to write file %q: %v", fnam, err)
		}
	}
	workMDFile := filepath.Join(workDir, mdFile)
	err = os.WriteFile(workMDFile, mdContent, 0666)
	if err != nil {
		ts.Fatalf("unable to write file %q: %v", workMDFile, err)
	}
	err = os.WriteFile(mdFile+".md", mdContent, 0666)
	if err != nil {
		ts.Fatalf("unable to write file %q: %v", mdFile+".md", err)
	}
}
