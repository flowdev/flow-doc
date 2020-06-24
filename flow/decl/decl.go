package decl

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"
	"unicode"

	"github.com/flowdev/ea-flow-doc/data"
	"github.com/flowdev/ea-flow-doc/flow/base"
)

// ParseFuncDecl parses a flow function (or method) declaration.
func ParseFuncDecl(decl *ast.FuncDecl, fset *token.FileSet, typesInfo *types.Info, flowDat *base.FlowData, errs []error,
) []error {

	flowDat.ComponentName, flowDat.InPort, errs = ParseFlowFuncName(decl.Name, fset, errs)
	log.Printf("DEBUG - componentName: %s, inPort: %v", flowDat.ComponentName, flowDat.InPort)

	flowDat.Inputs, errs = parseInputData(decl.Type.Params, fset, typesInfo, errs)
	for _, dat := range flowDat.Inputs {
		log.Printf("DEBUG - data: %v", dat)
	}

	var results []base.DataTyp
	results, flowDat.OutPorts, errs = parseFlowFuncResults(decl.Type.Results, fset, typesInfo, errs)
	flowDat.MainBranch.DataMap = base.AddDatasToMap(flowDat.MainBranch.DataMap, flowDat.Inputs)
	flowDat.MainBranch.DataMap = base.AddDatasToMap(flowDat.MainBranch.DataMap, results)
	for _, port := range flowDat.OutPorts {
		log.Printf("DEBUG - outPort: %v", port)
	}

	return errs
}

// ParseFlowFuncName parses a flow function name.
func ParseFlowFuncName(funcNameID *ast.Ident, fset *token.FileSet, errs []error,
) (componentName string, inPort base.Port, errs2 []error) {

	funcName := funcNameID.Name
	componentName = funcName
	inPort.Name = "in"
	inPort.IsImplicit = true

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
	inPort.Name = parts[1]
	inPort.Pos = funcNameID.Pos()
	inPort.IsImplicit = false

	if !unicode.IsLower([]rune(inPort.Name)[0]) {
		errs = append(errs, errors.New(
			fset.Position(funcNameID.Pos()).String()+
				" port names in flow function names must start with a lower case letter, got '"+
				inPort.Name+
				"' in: "+
				funcName,
		))
	}
	return componentName, inPort, errs
}

func parseInputData(params *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]base.DataTyp, []error) {

	if params == nil || len(params.List) == 0 {
		return nil, errs
	}

	var inputs []base.DataTyp

	inputs, errs = flowDataTypes(params, fset, typesInfo, errs)

	firstPlugin := -1
	for i, input := range inputs {
		if firstPlugin < 0 && isPlugin(input) {
			firstPlugin = i
		} else if firstPlugin >= 0 && !isPlugin(input) {
			errs = append(errs, errors.New(
				fset.Position(input.NamePos).String()+
					" flow plugins must all be at the end of the parameter list, found '"+
					input.Name+"' after plugin '"+inputs[firstPlugin].Name+"'",
			))
		}
	}

	return inputs, errs
}

func parseFlowFuncResults(funcResults *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]base.DataTyp, []base.Port, []error) {

	if funcResults == nil || len(funcResults.List) == 0 {
		return nil, nil, errs
	}

	portNames := 0
	defaultPort := base.Port{Name: "out", IsImplicit: true}
	lastIsError := false
	ports := []base.Port{}

	datas, _ := flowDataTypes(funcResults, fset, typesInfo, []error{})
	n := len(datas)

	if datas[n-1].Typ == "error" {
		lastIsError = true
	}

	for _, dat := range datas {
		if strings.HasPrefix(dat.Name, base.PortPrefix) && len(dat.Name) > len(base.PortPrefix) {
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
			ports = append(ports, base.Port{Name: portName(dat.Name), Pos: dat.NamePos})
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
		ports = append(ports, base.Port{Name: "error", IsError: true})
	}

	return datas, ports, errs
}

func flowDataTypes(fl *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]base.DataTyp, []error) {

	datas := make([]base.DataTyp, 0, 32)
	for _, field := range fl.List {
		flowDataType, err := data.Type(field.Type)
		if err != nil {
			errs = append(errs, errors.New(
				fset.Position(field.Type.Pos()).String()+
					" "+err.Error()+"; Go data type: "+
					base.TypeInfo(field.Type, typesInfo),
			))
			log.Printf("DEBUG - data type error: %s", // TODO: remove debug log
				fset.Position(field.Type.Pos()).String()+
					" "+err.Error()+"; Go data type: "+
					base.TypeInfo(field.Type, typesInfo))
		}
		for _, id := range field.Names {
			datas = append(datas, base.DataTyp{
				Name: id.Name, NamePos: id.NamePos,
				Typ: flowDataType, TypPos: field.Type.Pos(),
			})
		}
		if len(field.Names) == 0 {
			datas = append(datas, base.DataTyp{Typ: flowDataType})
		}
	}

	return datas, errs
}

func isPlugin(input base.DataTyp) bool {
	const prefixPlugin = "plugin"
	return strings.HasPrefix(input.Name, prefixPlugin) &&
		(len(input.Name) > len(prefixPlugin))
}

func portName(longName string) string {
	name := longName[len(base.PortPrefix):]
	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
