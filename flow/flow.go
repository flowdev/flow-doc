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
	name       string
	pos        token.Pos
	isImplicit bool
	isError    bool
}

type dataTyp struct {
	name    string
	namePos token.Pos
	typ     string
	typPos  token.Pos
}

type callStep struct {
	inputs        []string
	inPort        port
	componentName string
	outputs       []string
}

type returnStep struct {
	datas   []string
	outPort port
}

type step interface {
}

type branch struct {
	dataMap map[string]string
	steps   []step
	parent  *branch
}

// a flow has got a mainBranch an possibly many sub-branches.
// the main branch always starts with the first call expression of the function.
// sub-branches are created with if expressions.
// consequently a flow can't start with an if expression!
type flowData struct {
	inPort        port
	inputs        []dataTyp
	componentName string
	outPorts      []port
	mainBranch    *branch
}

func newBranch(parent *branch) *branch {
	return &branch{
		dataMap: make(map[string]string, 64),
		steps:   make([]step, 0, 64),
		parent:  parent,
	}
}
func newFlowData() *flowData {
	return &flowData{mainBranch: newBranch(nil)}
}

func parse(allFlowFuncs []find.PackageFuncs) ([]*flowData, []error) {
	var flowDatas []*flowData
	var allErrs []error

	for _, pkgFlowFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFlowFuncs.Funcs {
			var flowDat *flowData
			flowDat, allErrs = parseFlow(flowFunc, pkgFlowFuncs.Fset, pkgFlowFuncs.TypesInfo, allErrs)
			flowDatas = append(flowDatas, flowDat)
		}
	}

	for _, err := range allErrs {
		log.Printf("NOTICE - error: %v", err)
	}
	return flowDatas, allErrs
}

func parseFlow(
	flowFunc *ast.FuncDecl,
	fset *token.FileSet, typesInfo *types.Info,
	errs []error,
) (*flowData, []error) {
	flowDat := newFlowData()

	errs = parseFuncDecl(flowFunc, fset, typesInfo, flowDat, errs)
	errs = parseFuncBody(flowFunc.Body, fset, typesInfo, flowDat, flowDat.mainBranch, errs)

	return flowDat, errs
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
			log.Printf("DEBUG - data type error: %s", // TODO: remove debug log
				fset.Position(field.Type.Pos()).String()+
					" "+err.Error()+"; Go data type: "+
					typeInfo(field.Type, typesInfo))
		}
		for _, id := range field.Names {
			datas = append(datas, dataTyp{
				name: id.Name, namePos: id.NamePos,
				typ: flowDataType, typPos: field.Type.Pos(),
			})
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
			m = addDataToMap(dat.name, dat.typ, m)
		}
	}
	return m
}
func addDataToMap(name, typ string, m map[string]string) map[string]string {
	t, ok := m[name]
	if !ok || len(typ) > len(t) {
		m[name] = typ
	}
	return m
}

func isBoring(typ string) bool {
	switch typ { // simple builtin types are 'boring'
	case "bool", "byte", "complex64", "complex128", "float32", "float64",
		"int", "int8", "int16", "int32", "int64", "rune", "string",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "":
		return true
	default:
		return false
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
