package flow

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"github.com/flowdev/ea-flow-doc/find"
)

type port struct {
	name     string
	implicit bool
}

// Check validates the given flow functions.
// They have to satisfy the following rules:
func Check(allFlowFuncs []find.PackageFuncs) []error {
	var allErrs []error
	for _, pkgFlowFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFlowFuncs.Funcs {
			errs := checkFlow(flowFunc, pkgFlowFuncs.Fset, pkgFlowFuncs.TypesInfo)
			if len(errs) > 0 {
				allErrs = append(allErrs, errs...)
			}
		}
	}
	return allErrs
}

// Cases:
// - one input port [x]
//   - one output port
//   - error output port
//   - multiple output ports
// - multiple input ports: simple [x]
// - stateful components: no extra handling [x]
func checkFlow(flowFunc *ast.FuncDecl, fset *token.FileSet, typesInfo *types.Info) []error {
	var errs []error
	var componentName string
	var inPort port
	var outPorts []port

	componentName, inPort, errs = checkFlowFuncName(flowFunc.Name, fset, errs)

	outPorts, errs = checkFlowFuncResults(flowFunc.Type.Results, fset, typesInfo, errs)

	_, _ = componentName, inPort
	_ = outPorts

	return errs
}

func checkFlowFuncName(funcNameID *ast.Ident, fset *token.FileSet, errs []error,
) (componentName string, inPort port, errs2 []error) {

	funcName := funcNameID.Name
	componentName = funcName
	inPort.name = "in"
	inPort.implicit = true

	if !strings.Contains(funcName, "_") {
		return componentName, inPort, errs
	}
	parts := strings.Split(funcName, "_")
	if len(parts) != 2 {
		errs = append(errs, errors.New(
			fset.Position(funcNameID.Pos()).String()+
				" flow function names must contain at most one underscore ('_'), got: "+
				funcName,
		))
	}

	if parts[0] == "" {
		errs = append(errs, errors.New(
			fset.Position(funcNameID.Pos()).String()+
				" flow function names must contain a valid component name, got '' in: "+
				funcName,
		))
	}
	componentName = parts[0]

	if parts[1] == "" {
		errs = append(errs, errors.New(
			fset.Position(funcNameID.Pos()).String()+
				" flow function names with '_' must contain a valid port name, got '' in: "+
				funcName,
		))
	}
	inPort.name = parts[1]
	inPort.implicit = false

	if !unicode.IsLower([]rune(inPort.name)[0]) {
		errs = append(errs, errors.New(
			fset.Position(funcNameID.Pos()).String()+
				" port names in flow function names must start with a lower case letter, got '"+
				inPort.name+
				"' in: "+
				funcName,
		))
	}
	return componentName, inPort, errs
}

// Cases:
//   - one output port
//   - error output port
//   - multiple output ports
func checkFlowFuncResults(funcResults *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]port, []error) {

	if funcResults == nil || len(funcResults.List) == 0 {
		return nil, errs
	}

	const portPrefix = "port"
	portNamesCount := 0
	resultsCount := 0
	//defaultPort := port{name: "out", implicit: true}
	ports := []port{}

	for _, result := range funcResults.List {
		for _, id := range result.Names {
			name := id.Name
			resultsCount++
			if strings.HasPrefix(name, portPrefix) && len(name) > len(portPrefix) {
				portNamesCount++
			}
		}
		//fmt.Println("TYPEs:", typesInfo)
		fmt.Println("TYPE:", typesInfo.Types[result.Type].Type)
	}

	return ports, errs
}
