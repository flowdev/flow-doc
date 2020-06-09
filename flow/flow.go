package flow

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"log"
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

type flowData struct {
	inPort        port
	inputs        []dataTyp
	dataMap       map[string]string
	componentName string
	outPorts      []port
}

func parse(allFlowFuncs []find.PackageFuncs) ([]flowData, []error) {
	var flowDatas []flowData
	var allErrs []error

	for _, pkgFlowFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFlowFuncs.Funcs {
			var flowDat flowData
			flowDat, allErrs = parseFlow(flowFunc, pkgFlowFuncs.Fset, pkgFlowFuncs.TypesInfo, allErrs)
			flowDatas = append(flowDatas, flowDat)
		}
	}

	for _, err := range allErrs {
		log.Printf("NOTICE - error: %v", err)
	}
	return flowDatas, allErrs
}

// Cases:
// - one input port [x]
//   - one output port [x]
//   - error output port [x]
//   - multiple output ports [x]
// - multiple input ports: simple [x]
// - stateful components: no extra handling [x]
func parseFlow(flowFunc *ast.FuncDecl, fset *token.FileSet, typesInfo *types.Info, errs []error,
) (flowData, []error) {
	flowDat := flowData{}

	errs = parseFuncDecl(flowFunc, fset, typesInfo, &flowDat, errs)
	errs = parseFuncBody(flowFunc, fset, typesInfo, &flowDat, errs)

	return flowDat, errs
}

func parseFuncBody(decl *ast.FuncDecl, fset *token.FileSet, typesInfo *types.Info, flowDat *flowData, errs []error,
) []error {

	return errs
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

func addDatasToMap(m map[string]string, datas []dataTyp) map[string]string {
	for _, dat := range datas {
		if dat.name != "" {
			m[dat.name] = dat.typ
		}
	}
	return m
}

func dataTypToData(d dataTyp) string {
	switch d.typ { // don't report 'boring' data types
	case "bool", "byte", "complex64", "complex128", "float32", "float64",
		"int", "int8", "int16", "int32", "int64",
		"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr":
		return d.name
	default:
		return d.typ
	}
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
