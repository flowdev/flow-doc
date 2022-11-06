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
	const mdFile = "markdown.actual"
	workDir := ts.Getenv("WORK")

	svgContents, mdContent, err := draw.FromFlowData(
		BigTestFlowData,
		draw.FlowModeSVGLinks,
		800,
		true,
	)
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
	err = os.WriteFile("testFlow.md", mdContent, 0666)
	if err != nil {
		ts.Fatalf("unable to write file %q: %v", mdFile, err)
	}
}
