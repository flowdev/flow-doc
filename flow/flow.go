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
	"github.com/flowdev/ea-flow-doc/flow/base"
)

func parse(allFlowFuncs []find.PackageFuncs) ([]*base.FlowData, []error) {
	var flowDatas []*base.FlowData
	var allErrs []error

	for _, pkgFlowFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFlowFuncs.Funcs {
			var flowDat *base.FlowData
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
) (*base.FlowData, []error) {
	flowDat := base.NewFlowData()

	errs = parseFuncDecl(flowFunc, fset, typesInfo, flowDat, errs)
	errs = parseFuncBody(flowFunc.Body, fset, typesInfo, flowDat, flowDat.MainBranch, errs)

	return flowDat, errs
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
					typeInfo(field.Type, typesInfo),
			))
			log.Printf("DEBUG - data type error: %s", // TODO: remove debug log
				fset.Position(field.Type.Pos()).String()+
					" "+err.Error()+"; Go data type: "+
					typeInfo(field.Type, typesInfo))
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

func addDatasToMap(m map[string]string, datas []base.DataTyp) map[string]string {
	for _, dat := range datas {
		if dat.Name != "" {
			m = addDataToMap(dat.Name, dat.Typ, m)
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
	name := longName[len(base.PortPrefix):]
	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
