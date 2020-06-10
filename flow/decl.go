package flow

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"
	"unicode"
)

func parseFuncDecl(decl *ast.FuncDecl, fset *token.FileSet, typesInfo *types.Info, flowDat *flowData, errs []error,
) []error {

	flowDat.componentName, flowDat.inPort, errs = parseFlowFuncName(decl.Name, fset, errs)
	log.Printf("DEBUG - componentName: %s, inPort: %v", flowDat.componentName, flowDat.inPort)

	flowDat.inputs, errs = parseInputData(decl.Type.Params, fset, typesInfo, errs)
	for _, dat := range flowDat.inputs {
		log.Printf("DEBUG - data: %v", dat)
	}

	var results []dataTyp
	results, flowDat.outPorts, errs = parseFlowFuncResults(decl.Type.Results, fset, typesInfo, errs)
	flowDat.mainBranch.dataMap = addDatasToMap(flowDat.mainBranch.dataMap, flowDat.inputs)
	flowDat.mainBranch.dataMap = addDatasToMap(flowDat.mainBranch.dataMap, results)
	for _, port := range flowDat.outPorts {
		log.Printf("DEBUG - outPort: %v", port)
	}

	return errs
}

func parseFlowFuncName(funcNameID *ast.Ident, fset *token.FileSet, errs []error,
) (componentName string, inPort port, errs2 []error) {

	funcName := funcNameID.Name
	componentName = funcName
	inPort.name = "in"
	inPort.isImplicit = true

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
	inPort.pos = funcNameID.Pos()
	inPort.isImplicit = false

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

func parseInputData(params *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]dataTyp, []error) {

	if params == nil || len(params.List) == 0 {
		return nil, errs
	}

	var inputs []dataTyp

	inputs, errs = flowDataTypes(params, fset, typesInfo, errs)

	firstPlugin := -1
	for i, input := range inputs {
		if firstPlugin < 0 && isPlugin(input) {
			firstPlugin = i
		} else if firstPlugin >= 0 && !isPlugin(input) {
			errs = append(errs, errors.New(
				fset.Position(input.namePos).String()+
					" flow plugins must all be at the end of the parameter list, found '"+
					input.name+"' after plugin '"+inputs[firstPlugin].name+"'",
			))
		}
	}

	return inputs, errs
}

func parseFlowFuncResults(funcResults *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]dataTyp, []port, []error) {

	if funcResults == nil || len(funcResults.List) == 0 {
		return nil, nil, errs
	}

	portNames := 0
	defaultPort := port{name: "out", isImplicit: true}
	lastIsError := false
	ports := []port{}

	datas, _ := flowDataTypes(funcResults, fset, typesInfo, []error{})
	n := len(datas)

	if datas[n-1].typ == "error" {
		lastIsError = true
	}

	for _, dat := range datas {
		if strings.HasPrefix(dat.name, portPrefix) && len(dat.name) > len(portPrefix) {
			portNames++
		}
	}

	log.Printf("DEBUG - portNames: %d, n: %d, lastIsError: %t", portNames, n, lastIsError)
	for _, dat := range datas {
		log.Printf("DEBUG - data: %v", dat)
	}

	if portNames == n || (portNames == n-1 && lastIsError) {
		for i, dat := range datas {
			if i == n-1 && lastIsError {
				break
			}
			ports = append(ports, port{name: portName(dat.name), pos: dat.namePos})
		}
	} else if n > 1 || (n == 1 && !lastIsError) {
		ports = append(ports, defaultPort)
		if portNames > 0 {
			position := ""
			if len(funcResults.List[0].Names) > 0 {
				position = fset.Position(funcResults.List[0].Names[0].NamePos).String()
			} else {
				position = fset.Position(funcResults.List[0].Type.Pos()).String()
			}
			log.Printf("WARNING - found only %d port names at: %s", portNames, position)
		}
	}

	if lastIsError {
		ports = append(ports, port{name: "error", isError: true})
	}

	return datas, ports, errs
}

func isPlugin(input dataTyp) bool {
	const prefixPlugin = "plugin"
	return strings.HasPrefix(input.name, prefixPlugin) &&
		(len(input.name) > len(prefixPlugin))
}
