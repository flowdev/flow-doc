package data_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/flowdev/ea-flow-doc/data"
	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/parse"
)

func TestType(t *testing.T) {
	specs := map[string]struct {
		expectedParams  string
		expectedResults string
		expectedErrors  int
	}{
		"simpleFunc": {
			expectedParams:  `[]`,
			expectedResults: `[]`,
			expectedErrors:  0,
		},
		"doAccountingMagic": {
			expectedParams:  `["string" "string"]`,
			expectedResults: `["string"]`,
			expectedErrors:  0,
		},
		"doSpecialAccountingMagic": {
			expectedParams:  `["string" "string" "BankAccount"]`,
			expectedResults: `["string" "SpecialBankAccount" "error"]`,
			expectedErrors:  0,
		},
		"funcWithEllipsis_in2": {
			expectedParams:  `["int" "list(string)" "list(string)" "list(bool)"]`,
			expectedResults: `["list(string)" "error"]`,
			expectedErrors:  0,
		},
		"funcWithErrorOnly": {
			expectedParams:  `[]`,
			expectedResults: `["error"]`,
			expectedErrors:  0,
		},
		"funcWithoutResults": {
			expectedParams:  `["tool.Data" "map(string, SpecialBankAccount)" "list(map(string, map(int, list(tool.Data))))" "string"]`,
			expectedResults: `[]`,
			expectedErrors:  0,
		},
		"funcWithTooComplexData": {
			expectedParams:  `["string" "" "" "int"]`,
			expectedResults: `[]`,
			expectedErrors:  2,
		},
	}
	root := mustAbs(filepath.Join("testdata", "typ"))
	pkgs, err := parse.Dir(root, true)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	allFlowFuncs := find.FlowFuncs(pkgs)

	for _, pkgFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFuncs.Funcs {
			name := flowFunc.Name.Name
			t.Run(name, func(t *testing.T) {
				actualErrors := 0
				actualParams := []string{}
				actualResults := []string{}

				for _, field := range flowFunc.Type.Params.List {
					typ, err := data.Type(field.Type)
					if err != nil {
						actualErrors++
					}
					actualParams = append(actualParams, typ)
				}
				if flowFunc.Type.Results != nil {
					for _, field := range flowFunc.Type.Results.List {
						typ, err := data.Type(field.Type)
						if err != nil {
							actualErrors++
						}
						actualResults = append(actualResults, typ)
					}
				}

				spec, ok := specs[name]
				if !ok {
					t.Fatalf("spec for function %q is missing", name)
				}
				actualParamsStr := fmt.Sprintf("%q", actualParams)
				actualResultsStr := fmt.Sprintf("%q", actualResults)

				if actualParamsStr != spec.expectedParams {
					t.Errorf("expected params %s, got: %s", spec.expectedParams, actualParamsStr)
				}
				if actualResultsStr != spec.expectedResults {
					t.Errorf("expected results %s, got: %s", spec.expectedResults, actualResultsStr)
				}
				if actualErrors != spec.expectedErrors {
					t.Errorf("expected errors %d, got: %d", spec.expectedErrors, actualErrors)
				}
			})
		}
	}
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
