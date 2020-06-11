package flow

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"unicode"

	"github.com/flowdev/ea-flow-doc/data"
	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/x/reflect"
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

type call struct {
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

func newBranch() *branch {
	return &branch{dataMap: make(map[string]string, 64), steps: make([]step, 0, 64)}
}
func newFlowData() *flowData {
	return &flowData{mainBranch: newBranch()}
}

// Parse is farsing flows.
func Parse(allFlowFuncs []find.PackageFuncs) ([]*flowData, []error) {
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

// Cases:
// - one input port [x]
//   - one output port [x]
//   - error output port [x]
//   - multiple output ports [x]
// - multiple input ports: simple [x]
// - stateful components: no extra handling [x]
func parseFlow(flowFunc *ast.FuncDecl, fset *token.FileSet, typesInfo *types.Info, errs []error,
) (*flowData, []error) {
	flowDat := newFlowData()

	errs = parseFuncDecl(flowFunc, fset, typesInfo, flowDat, errs)
	errs = parseFuncBody(flowFunc.Body, fset, typesInfo, flowDat, errs)

	return flowDat, errs
}

// BODY: -----------------------------

func parseFuncBody(body *ast.BlockStmt, fset *token.FileSet, typesInfo *types.Info, flowDat *flowData, errs []error,
) []error {

	for _, stmt := range body.List {
		errs = parseFuncStmt(stmt, fset, typesInfo, flowDat, flowDat.mainBranch, errs)
	}
	return errs
}

func parseFuncStmt(stmt ast.Stmt, fset *token.FileSet, typesInfo *types.Info, flowDat *flowData, branch *branch, errs []error,
) []error {

	if reflect.IsNilInterfaceOrPointer(stmt) {
		return errs
	}

	switch s := stmt.(type) {
	case *ast.DeclStmt:
		errs = parseDecl(s.Decl, fset, typesInfo, branch, errs)
	case *ast.ExprStmt:
		// TODO: only allow CallExpr!
	case *ast.AssignStmt:
		// TODO: Rhs: allow only calls?
	case *ast.ReturnStmt:
		// TODO: check Results: is error given? What out port is used?
	case *ast.IfStmt:
		// TODO: for error handling
	case *ast.ForStmt,
		*ast.RangeStmt,
		*ast.BlockStmt,
		*ast.SwitchStmt,
		*ast.TypeSwitchStmt,
		*ast.CaseClause,
		*ast.SelectStmt,
		*ast.CommClause,
		*ast.SendStmt,
		*ast.BranchStmt,
		*ast.GoStmt,
		*ast.LabeledStmt,
		*ast.DeferStmt,
		*ast.IncDecStmt:

		errs = append(errs, errors.New(
			fset.Position(stmt.Pos()).String()+
				" unsupported statement in flow, allowed are: "+
				"variable declaration, assignment, function calls, return and 'if <port>!=nil'",
		))
	case *ast.EmptyStmt,
		nil:
		// nothing to do
	default:
		errs = append(errs, errors.New(
			fset.Position(stmt.Pos()).String()+
				fmt.Sprintf("don't know how to handle unknown statement in flow: %T", s),
		))
	}
	return errs
}

func parseDecl(decl ast.Decl, fset *token.FileSet, typesInfo *types.Info, branch *branch, errs []error,
) []error {

	if reflect.IsNilInterfaceOrPointer(decl) {
		return errs
	}

	switch d := decl.(type) {
	case *ast.FuncDecl:
		errs = append(errs, errors.New(
			fset.Position(decl.Pos()).String()+
				" function declarations aren't supported in flows, allowed are: "+
				"variable declaration, assignment, function calls, return and 'if <port>!=nil'",
		))
	case *ast.GenDecl:
		errs = parseGenDecl(d, fset, typesInfo, branch, errs)
	default:
		errs = append(errs, errors.New(
			fset.Position(decl.Pos()).String()+
				fmt.Sprintf("don't know how to handle unknown declaration in flow: %T", d),
		))
	}

	return errs
}

func parseGenDecl(decl *ast.GenDecl, fset *token.FileSet, typesInfo *types.Info, branch *branch, errs []error,
) []error {

	if reflect.IsNilInterfaceOrPointer(decl) {
		return errs
	}

	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			errs = append(errs, errors.New(
				fset.Position(spec.Pos()).String()+
					" function declarations aren't supported in flows, allowed are: "+
					"variable declaration, assignment, function calls, return and 'if <port>!=nil'",
			))
		case *ast.ValueSpec:
			var typ string
			var err error
			if s.Type != nil {
				if typ, err = data.Type(s.Type); err != nil {
					errs = append(errs, errors.New(
						fset.Position(s.Type.Pos()).String()+
							" "+err.Error()+"; Go data type: "+
							typeInfo(s.Type, typesInfo),
					))
				}
			}
			for _, n := range s.Names {
				branch.dataMap[n.Name] = typ
			}
		}
		//default: import specs are ignored
	}

	return errs
}

// BODY: -----------------------------

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
			m[dat.name] = dat.typ
		}
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
