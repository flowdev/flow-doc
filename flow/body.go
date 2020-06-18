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
