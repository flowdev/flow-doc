package flow

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"
	"unicode"

	"github.com/flowdev/ea-flow-doc/data"
	"github.com/flowdev/ea-flow-doc/find"
)

const portPrefix = "port"

type port struct {
	name     string
	implicit bool
}

type dataTyp struct {
	name string
	typ  string
}

// Check validates the given flow functions.
// They have to satisfy the following rules:
func Check(allFlowFuncs []find.PackageFuncs) []error {
	var allErrs []error
	for _, pkgFlowFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFlowFuncs.Funcs {
			allErrs = checkFlow(flowFunc, pkgFlowFuncs.Fset, pkgFlowFuncs.TypesInfo, allErrs)
		}
	}

	for _, err := range allErrs {
		log.Printf("NOTICE - error: %v", err)
	}
	return allErrs
}

// Cases:
// - one input port [x]
//   - one output port [x]
//   - error output port [x]
//   - multiple output ports [x]
// - multiple input ports: simple [x]
// - stateful components: no extra handling [x]
func checkFlow(flowFunc *ast.FuncDecl, fset *token.FileSet, typesInfo *types.Info, errs []error) []error {
	var componentName string
	var inPort port
	var datas []dataTyp
	var outPorts []port

	componentName, inPort, errs = checkFlowFuncName(flowFunc.Name, fset, errs)

	log.Printf("DEBUG - componentName: %s, inPort: %v", componentName, inPort)
	datas, errs = checkInputData(flowFunc.Type.Params, fset, typesInfo, errs)
	for _, dat := range datas {
		log.Printf("DEBUG - data: %v", dat)
	}

	outPorts, errs = checkFlowFuncResults(flowFunc.Type.Results, fset, typesInfo, errs)
	for _, port := range outPorts {
		log.Printf("DEBUG - outPort: %v", port)
	}

	_, _ = componentName, inPort
	_ = datas
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

func checkInputData(params *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]dataTyp, []error) {

	if params == nil || len(params.List) == 0 {
		return nil, errs
	}

	var datas []dataTyp

	datas, errs = flowDataTypes(params, fset, typesInfo, errs)
	return datas, errs
}

func checkFlowFuncResults(funcResults *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]port, []error) {

	if funcResults == nil || len(funcResults.List) == 0 {
		return nil, errs
	}

	portNames := 0
	defaultPort := port{name: "out", implicit: true}
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
			ports = append(ports, port{name: portName(dat.name)})
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
		ports = append(ports, port{name: "error"})
	}

	return ports, errs
}

func flowDataTypes(fl *ast.FieldList, fset *token.FileSet, typesInfo *types.Info, errs []error,
) ([]dataTyp, []error) {

	datas := make([]dataTyp, 0, 32)
	for _, field := range fl.List {
		flowDataType, err := data.Type(field.Type)
		if err != nil {
			errs = append(errs, errors.New(
				fset.Position(field.Type.Pos()).String()+
					" "+err.Error()+"; Go data type: "+
					typeInfo(field.Type, typesInfo),
			))
			log.Printf("DEBUG - data type error: %s",
				fset.Position(field.Type.Pos()).String()+
					" "+err.Error()+"; Go data type: "+
					typeInfo(field.Type, typesInfo))
		}
		for _, id := range field.Names {
			datas = append(datas, dataTyp{name: id.Name, typ: flowDataType})
		}
		if len(field.Names) == 0 {
			datas = append(datas, dataTyp{typ: flowDataType})
		}
	}

	return datas, errs
}

func typeInfo(typ ast.Expr, typesInfo *types.Info) string {
	if typesInfo.Types[typ].Type == nil {
		return "<types.Info not filled properly>"
	}
	return typesInfo.Types[typ].Type.String()
}

func portName(longName string) string {
	name := longName[len(portPrefix):]
	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
