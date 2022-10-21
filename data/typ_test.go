package data_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/flowdev/ea-flow-doc/data"
	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/parse"
)

const actualResultsKey = "actualResults"

type actualResult struct {
	name    string
	params  string
	results string
	errors  int
}

func setup(e *testscript.Env) error {
	pkgs, err := parse.Dir(e.WorkDir, true)
	if err != nil {
		return fmt.Errorf("received unexpected parse error: %w", err)
	}

	actualResults := make(map[string]actualResult)
	allFlowFuncs := find.FlowFuncs(pkgs)

	for _, pkgFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFuncs.Funcs {
			errors := 0
			actualParams := []string{}
			returns := []string{}

			for _, field := range flowFunc.Type.Params.List {
				typ, err := data.Type(field.Type)
				if err != nil {
					errors++
				}
				actualParams = append(actualParams, typ)
			}
			if flowFunc.Type.Results != nil {
				for _, field := range flowFunc.Type.Results.List {
					typ, err := data.Type(field.Type)
					if err != nil {
						errors++
					}
					returns = append(returns, typ)
				}
			}

			result := actualResult{
				name:    flowFunc.Name.Name,
				params:  fmt.Sprintf("%q", actualParams),
				results: fmt.Sprintf("%q", returns),
				errors:  errors,
			}
			actualResults[result.name] = result
		}
	}
	e.Values[actualResultsKey] = actualResults
	return nil
}

func testTypeFunc(
	name string,
	expectedParams, expectedResults string,
	expectedErrors int,
	actualResults map[string]actualResult,
) error {
	act, ok := actualResults[name]
	if !ok {
		return fmt.Errorf("actual values for function %q are missing", name)
	}

	errMsg := &strings.Builder{}

	if act.params != expectedParams {
		errMsg.WriteString(fmt.Sprintf("expected params %s, got: %s\n", expectedParams, act.params))
	}
	if act.results != expectedResults {
		errMsg.WriteString(fmt.Sprintf("expected results %s, got: %s\n", expectedResults, act.results))
	}
	if act.errors != expectedErrors {
		errMsg.WriteString(fmt.Sprintf("expected errors %d, got: %d\n", expectedErrors, act.errors))
	}

	if errMsg.Len() != 0 {
		return errors.New(errMsg.String())
	}
	return nil
}

func TestType(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:   "testdata",
		Setup: setup,
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"expectTypeFunc": func(ts *testscript.TestScript, _ bool, args []string) {
				if len(args) != 4 {
					ts.Fatalf("expected 4 arguments ("+
						"name, expectedParams, expectedResults, expectedErrors"+
						") but got: %q", args)
				}
				errCount, err := strconv.Atoi(args[3])
				if err != nil {
					ts.Fatalf("expected argument 'expectedErrors' to be an integer number but got %q: %v",
						args[3], err)
				}
				testTypeFunc(args[0], args[1], args[2], errCount, ts.Value(actualResultsKey).(map[string]actualResult))
			},
		},
		TestWork: false,
	})
}
