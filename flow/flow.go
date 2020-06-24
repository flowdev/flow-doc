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

	errs = decl.ParseFuncDecl(flowFunc, fset, typesInfo, flowDat, errs)
	errs = body.ParseFuncBody(flowFunc.Body, fset, typesInfo, flowDat, flowDat.MainBranch, errs)

	return flowDat, errs
}
