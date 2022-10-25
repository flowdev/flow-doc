package flow

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"github.com/flowdev/ea-flow-doc/find"
	"github.com/flowdev/ea-flow-doc/flow/base"
	"github.com/flowdev/ea-flow-doc/flow/body"
	"github.com/flowdev/ea-flow-doc/flow/decl"
)

func parseAll(allFlowFuncs []find.PackageFuncs) ([]*base.FlowData, []error) {
	var flowDatas []*base.FlowData
	var allErrs []error

	for _, pkgFlowFuncs := range allFlowFuncs {
		for _, flowFunc := range pkgFlowFuncs.Funcs {
			flowDat, errs := parseFlowFunc(flowFunc, pkgFlowFuncs.Fset, pkgFlowFuncs.TypesInfo)
			flowDatas = append(flowDatas, flowDat)
			allErrs = append(allErrs, errs...)
		}
	}

	for _, err := range allErrs {
		log.Printf("NOTICE - error: %v", err)
	}
	return flowDatas, allErrs
}

func parseFlowFunc(
	flowFunc *ast.FuncDecl,
	fset *token.FileSet, typesInfo *types.Info,
) (*base.FlowData, []error) {
	errs := make([]error, 0, 32)
	flowDat := base.NewFlowData()

	errs = decl.ParseFuncDecl(flowFunc, fset, typesInfo, flowDat, errs)
	errs = body.ParseFuncBody(flowFunc.Body, fset, typesInfo, flowDat, flowDat.MainBranch, errs)

	return flowDat, errs
}
