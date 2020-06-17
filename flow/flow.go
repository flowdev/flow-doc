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
	errs = parseFuncBody(flowFunc.Body, fset, typesInfo, flowDat, errs)

	return flowDat, errs
}

// BODY: -----------------------------

func parseFuncBody(
	body *ast.BlockStmt,
	fset *token.FileSet, typesInfo *types.Info,
	flowDat *flowData,
	errs []error,
) []error {

	for _, stmt := range body.List {
		errs = parseFuncStmt(stmt, fset, typesInfo, flowDat, flowDat.mainBranch, errs)
	}
	return errs
}

func parseFuncStmt(
	stmt ast.Stmt,
	fset *token.FileSet, typesInfo *types.Info,
	flowDat *flowData, branch *branch,
	errs []error,
) []error {

	if reflect.IsNilInterfaceOrPointer(stmt) {
		return errs
	}

	switch s := stmt.(type) {
	case *ast.DeclStmt:
		errs = parseDecl(s.Decl, fset, typesInfo, branch, errs)
	case *ast.ExprStmt:
		var call *callStep
		call, errs = parseCall(s.X, false, fset, errs)
		if call != nil {
			branch.steps = append(branch.steps, call)
		}
	case *ast.AssignStmt:
		errs = parseAssignLhs(s.Lhs, fset, branch, errs)
		if len(s.Rhs) == 1 {
			var call *callStep
			call, errs = parseCall(s.Rhs[0], true, fset, errs)
			if call != nil {
				branch.steps = append(branch.steps, call)
			}
		} else {
			errs = parseAssignRhs(s.Rhs, fset, errs)
		}
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
					" type declarations aren't supported in flows, allowed are: "+
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

func parseCall(
	expr ast.Expr, allowLiteral bool,
	fset *token.FileSet,
	errs []error,
) (*callStep, []error) {

	if reflect.IsNilInterfaceOrPointer(expr) {
		pos := "<unknown position>"
		if expr != nil {
			pos = fset.Position(expr.Pos()).String()
		}
		errs = append(errs, errors.New(pos+
			fmt.Sprintf(" missing call expression in flow"),
		))
		return nil, errs
	}

	var call *callStep

	switch e := expr.(type) {
	case *ast.CallExpr:
		call = &callStep{}
		// check function name:
		var funcNameID *ast.Ident
		funcNameID, errs = getFunctionNameID(e.Fun, fset, errs)
		if funcNameID != nil {
			call.componentName, call.inPort, errs = parseFlowFuncName(funcNameID, fset, errs)
		}
		call.inputs, errs = getFunctionArguments(e.Args, fset, errs)
	case *ast.BasicLit:
		if !allowLiteral {
			errs = append(errs, errors.New(
				fset.Position(expr.Pos()).String()+
					fmt.Sprintf("don't know how to handle literal at this position in flow: %T", e),
			))
		}
	case *ast.Ident:
		if !allowLiteral {
			errs = append(errs, errors.New(
				fset.Position(expr.Pos()).String()+
					fmt.Sprintf("don't know how to handle identifier at this position in flow: %T", e),
			))
		}
	case nil:
		// should be very rare
		log.Printf("DEBUG - %s nil expression found", fset.Position(expr.Pos()).String())
		errs = append(errs, errors.New(
			fset.Position(expr.Pos()).String()+
				fmt.Sprintf(" nil expression found in flow"),
		))
	default:
		errs = append(errs, errors.New(
			fset.Position(expr.Pos()).String()+
				fmt.Sprintf("don't know how to handle unknown expression in flow: %T", e),
		))
	}

	return call, errs
}

func getFunctionNameID(expr ast.Expr, fset *token.FileSet, errs []error,
) (*ast.Ident, []error) {

	if reflect.IsNilInterfaceOrPointer(expr) {
		pos := "<unknown position>"
		if expr != nil {
			pos = fset.Position(expr.Pos()).String()
		}
		errs = append(errs, errors.New(pos+
			fmt.Sprintf("missing function name in call expression in flow"),
		))
		return nil, errs
	}

	switch e := expr.(type) {
	case *ast.Ident:
		return e, errs
	default:
		errs = append(errs, errors.New(
			fset.Position(expr.Pos()).String()+
				fmt.Sprintf(
					"can't find function name in call expression in flow, got: %T", e,
				),
		))
	}
	return nil, errs
}

func getFunctionArguments(args []ast.Expr, fset *token.FileSet, errs []error,
) ([]string, []error) {

	strArgs := make([]string, len(args))
	for i, arg := range args {
		strArgs[i], errs = parseIdent(arg, fset, "function argument in call expression", errs)
	}
	return strArgs, errs
}

func parseIdent(expr ast.Expr, fset *token.FileSet, errMsg string, errs []error,
) (string, []error) {

	if reflect.IsNilInterfaceOrPointer(expr) {
		pos := "<unknown position>"
		if expr != nil {
			pos = fset.Position(expr.Pos()).String()
		}
		errs = append(errs, errors.New(pos+fmt.Sprintf("missing %s in flow", errMsg)))
		return "<error>", errs
	}

	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name, errs
	default:
		errs = append(errs, errors.New(
			fset.Position(expr.Pos()).String()+
				fmt.Sprintf("can't find %s in flow, got: %T", errMsg, e),
		))
		return "<error>", errs
	}
}

func parseAssignLhs(exprs []ast.Expr, fset *token.FileSet, branch *branch, errs []error,
) []error {
	for _, expr := range exprs {
		id := ""
		id, errs = parseIdent(expr, fset, "identifier in assignment", errs)
		if id != "" {
			branch.dataMap = addDataToMap(id, "", branch.dataMap)
		}
	}
	return errs
}

func parseAssignRhs(exprs []ast.Expr, fset *token.FileSet, errs []error) []error {
	for _, expr := range exprs {
		errs = parseSimpleExpression(expr, fset, errs)
	}
	return errs
}

func parseSimpleExpression(
	expr ast.Expr,
	fset *token.FileSet,
	errs []error,
) []error {

	if reflect.IsNilInterfaceOrPointer(expr) {
		pos := "<unknown position>"
		if expr != nil {
			pos = fset.Position(expr.Pos()).String()
		}
		errs = append(errs, errors.New(pos+
			fmt.Sprintf(" missing simple expression in right hand side of assignmnet in flow"),
		))
		return errs
	}

	switch e := expr.(type) {
	case *ast.BasicLit:
		// all good
	case *ast.Ident:
		// all good
	case nil:
		// should be very rare
		log.Printf("DEBUG - %s nil expression found", fset.Position(expr.Pos()).String())
		errs = append(errs, errors.New(
			fset.Position(expr.Pos()).String()+
				fmt.Sprintf(" nil expression found in flow"),
		))
	default:
		errs = append(errs, errors.New(
			fset.Position(expr.Pos()).String()+
				fmt.Sprintf(
					"don't know how to handle unknown expression in right hand side of assignment in flow: %T",
					e,
				),
		))
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
