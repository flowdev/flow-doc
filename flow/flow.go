package flow

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/flowdev/ea-flow-doc/find"
)

// Check validates the given flow functions.
// They have to satisfy the following rules:
func Check(allFlowFuncs []find.PackageFuncs) []error {
	var allErrs []error
	for _, pkgFlowFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFlowFuncs.Funcs {
			errs := checkFlow(flowFunc, pkgFlowFuncs.Fset)
			if len(errs) > 0 {
				allErrs = append(allErrs, errs...)
			}
		}
	}
	return allErrs
}

// Cases:
// - one input port
//   - one output port
//   - error output port
//   - multiple output ports
// - multiple input ports
// - stateful components: no extra handling
func checkFlow(flowFunc *ast.FuncDecl, fset *token.FileSet) []error {
	var errs []error

	funcName := flowFunc.Name.Name
	componentName := funcName
	portName := "in"
	implicitPort := true

	if strings.Contains(funcName, "_") {
		parts := strings.Split(funcName, "_")
		if len(parts) != 2 {
			errs = append(errs, errors.New(
				fset.Position(flowFunc.Name.Pos()).String()+
					" flow function names must contain at most one underscore ('_'), got: "+
					funcName,
			))
		}

		if parts[0] == "" {
			errs = append(errs, errors.New(
				fset.Position(flowFunc.Name.Pos()).String()+
					" flow function names must contain a valid component name, got '' in: "+
					funcName,
			))
		}
		componentName = parts[0]

		if parts[1] == "" {
			errs = append(errs, errors.New(
				fset.Position(flowFunc.Name.Pos()).String()+
					" flow function names with '_' must contain a valid port name, got '' in: "+
					funcName,
			))
		}
		portName = parts[1]
		implicitPort = false

		if !unicode.IsLower([]rune(portName)[0]) {
			errs = append(errs, errors.New(
				fset.Position(flowFunc.Name.Pos()).String()+
					" port names in flow function names must start with a lower case letter, got '"+
					portName+
					"' in: "+
					funcName,
			))
		}
	}
	return errs
}
