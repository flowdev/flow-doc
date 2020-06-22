package flow

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"github.com/flowdev/ea-flow-doc/data"
	"github.com/flowdev/ea-flow-doc/x/reflect"
)

const identNameError = "<error>"

type identType int

const (
	identTypeStrict identType = iota
	identTypeOrUnderscore
	identTypeOrNil
	identTypeOnlyNil
)

func parseFuncBody(
	body *ast.BlockStmt,
	fset *token.FileSet, typesInfo *types.Info,
	flowDat *flowData, branch *branch,
	errs []error,
) []error {

	for _, stmt := range body.List {
		branch, errs = parseFuncStmt(stmt, fset, typesInfo, flowDat, branch, errs)
	}
	return errs
}

func parseFuncStmt(
	stmt ast.Stmt,
	fset *token.FileSet, typesInfo *types.Info,
	flowDat *flowData, branch *branch,
	errs []error,
) (*branch, []error) {

	if reflect.IsNilInterfaceOrPointer(stmt) {
		return branch, errs
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
		errs = parseReturn(s, fset, flowDat, branch, errs)
		if branch.parent != nil {
			branch = branch.parent
		}
	case *ast.IfStmt:
		branch, errs = parseIf(s, fset, typesInfo, flowDat, branch, errs)
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
	return branch, errs
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
		strArgs[i], errs = parseIdent(arg, identTypeOrNil, fset, "function argument in call expression", errs)
	}
	return strArgs, errs
}

func parseIdent(expr ast.Expr, idTyp identType, fset *token.FileSet, errMsg string, errs []error,
) (string, []error) {

	if reflect.IsNilInterfaceOrPointer(expr) {
		pos := "<unknown position>"
		if expr != nil {
			pos = fset.Position(expr.Pos()).String()
		}
		errs = append(errs, errors.New(pos+fmt.Sprintf("missing %s in flow", errMsg)))
		return identNameError, errs
	}

	switch e := expr.(type) {
	case *ast.Ident:
		if idTyp == identTypeOnlyNil {
			errs = append(errs, errors.New(
				fset.Position(expr.Pos()).String()+
					fmt.Sprintf("can't find %s in flow, got: %q", errMsg, e.Name),
			))
			return identNameError, errs
		}
		return e.Name, errs
		// TODO: handle nil case
		// TODO: handle _ case
	default:
		errs = append(errs, errors.New(
			fset.Position(expr.Pos()).String()+
				fmt.Sprintf("can't find %s in flow, got: %T", errMsg, e),
		))
		return identNameError, errs
	}
}

func parseAssignLhs(exprs []ast.Expr, fset *token.FileSet, branch *branch, errs []error,
) []error {
	for _, expr := range exprs {
		id := ""
		id, errs = parseIdent(expr, identTypeOrUnderscore, fset, "identifier in assignment", errs)
		if id != identNameError {
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

func parseReturn(
	ret *ast.ReturnStmt,
	fset *token.FileSet,
	flowDat *flowData, branch *branch,
	errs []error,
) []error {
	ops := flowDat.outPorts
	opsN := len(ops)
	resM := len(ret.Results) - 1

	if opsN == 0 { // no output at all
		// nothing to do
	} else if opsN == 1 && ops[0].isImplicit { // only 'out'
		errs = parseImplicitOutPort(
			ret.Results,
			ops[0], flowDat.mainBranch.dataMap,
			fset, branch,
			errs,
		)
	} else if opsN == 2 && ops[0].isImplicit && ops[1].isError { // 'out' && 'error'
		if len(ret.Results) == 0 {
			errs = append(errs, errors.New(fset.Position(ret.Return).String()+
				fmt.Sprintf(" missing value in return statement in flow"),
			))
			return errs
		}
		// check error first:
		done := false
		done, errs = parseExplicitPort(
			ret.Results[resM], false,
			ops[1], flowDat.mainBranch.dataMap,
			fset, branch,
			errs,
		)
		if done {
			return errs
		}

		errs = parseImplicitOutPort(
			ret.Results[:resM],
			ops[0], flowDat.mainBranch.dataMap,
			fset, branch,
			errs,
		)
	} else { // explicit ports (including error)
		if len(ret.Results) == 0 {
			errs = append(errs, errors.New(fset.Position(ret.Return).String()+
				fmt.Sprintf(" missing value in return statement in flow"),
			))
			return errs
		}

		if opsN != resM+1 {
			errs = append(errs, errors.New(fset.Position(ret.Return).String()+
				fmt.Sprintf(" %d return values don't match %d output ports", resM+1, opsN),
			))
			return errs
		}
		found := false
		for i := 0; i <= resM; i++ {
			found, errs = parseExplicitPort(
				ret.Results[i], found,
				ops[i], flowDat.mainBranch.dataMap,
				fset, branch,
				errs,
			)
		}
		if found {
			return errs
		}
		errs = append(errs, errors.New(fset.Position(ret.Return).String()+
			fmt.Sprintf(" no port of %d possible ports selected in return statement", opsN),
		))
		return errs
	}
	return errs
}

func parseExplicitPort(
	result ast.Expr, found bool,
	op port, globalData map[string]string,
	fset *token.FileSet,
	branch *branch,
	errs []error,
) (done bool, errs2 []error) {
	name := ""
	name, errs = parseIdent(result, identTypeOrNil, fset, "name in return statement", errs)
	if name != identNameError {
		if found {
			errs = append(errs, errors.New(fset.Position(result.Pos()).String()+
				fmt.Sprintf(
					" found value %q for port %q even though another port has been sent to already",
					name, op.name,
				),
			))
			return true, errs
		}
		branch.steps = append(branch.steps,
			&returnStep{
				datas:   []string{dataForName(name, branch.dataMap, globalData)},
				outPort: op,
			})
		return true, errs
	}
	return found, errs
}

func parseImplicitOutPort(
	results []ast.Expr,
	op port, globalData map[string]string,
	fset *token.FileSet,
	branch *branch,
	errs []error,
) []error {

	name := ""
	rs := &returnStep{datas: make([]string, 0, len(results)), outPort: op}
	for _, result := range results {
		name, errs = parseIdent(result, identTypeOrNil, fset, "name in return statement", errs)
		if name != identNameError {
			rs.datas = append(rs.datas, dataForName(name, branch.dataMap, globalData))
		}
	}
	branch.steps = append(branch.steps, rs)
	return errs
}

func parseIf(
	ifs *ast.IfStmt,
	fset *token.FileSet, typesInfo *types.Info,
	flowDat *flowData, branch *branch,
	errs []error,
) (*branch, []error) {
	if ifs.Else != nil {
		errs = append(errs, errors.New(fset.Position(ifs.Else.Pos()).String()+
			" else branch of 'if' statement isn't allowed in flows"),
		)
	}
	errs = parseIfCond(ifs.Cond, fset, errs)
	branch = newBranch(branch)
	errs = parseFuncBody(ifs.Body, fset, typesInfo, flowDat, branch, errs)
	return branch, errs
}

func parseIfCond(
	cond ast.Expr,
	fset *token.FileSet,
	errs []error,
) []error {

	if reflect.IsNilInterfaceOrPointer(cond) {
		pos := "<unknown position>"
		if cond != nil {
			pos = fset.Position(cond.Pos()).String()
		}
		errs = append(errs, errors.New(pos+
			fmt.Sprintf(" missing condition in if statement in flow"),
		))
		return errs
	}

	switch e := cond.(type) {
	case *ast.BinaryExpr:
		// ident != nil
		parseIfCondition(e, fset, errs)
	case nil:
		// should be very rare
		log.Printf("DEBUG - %s nil expression found", fset.Position(cond.Pos()).String())
		errs = append(errs, errors.New(
			fset.Position(cond.Pos()).String()+
				fmt.Sprintf(" nil expression found in flow"),
		))
	default:
		errs = append(errs, errors.New(
			fset.Position(cond.Pos()).String()+
				fmt.Sprintf(
					"don't know how to handle unknown expression in if condition in flow: %T",
					e,
				),
		))
	}

	return errs
}

func parseIfCondition(
	be *ast.BinaryExpr,
	fset *token.FileSet,
	errs []error,
) []error {

	if be.Op != token.NEQ {
		errs = append(errs, errors.New(
			fset.Position(be.OpPos).String()+
				fmt.Sprintf(
					"only \"!=\" allowed as operator in if condition in flows, got: %q",
					be.Op.String(),
				),
		))
	}
	_, errs = parseIdent(be.X, identTypeStrict, fset, "name in if condition", errs)
	_, errs = parseIdent(be.Y, identTypeOnlyNil, fset, "nil in if condition", errs)

	return errs
}

func dataForName(name string, localData, globalData map[string]string) string {
	if d, ok := localData[name]; ok {
		if d == "" {
			return name
		}
		return d
	}
	if d, ok := globalData[name]; ok {
		if d == "" {
			return name
		}
		return d
	}
	return name
}
