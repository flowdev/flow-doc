package parse

import (
	"go/ast"

	"github.com/flowdev/ea-flow-doc/x/pkgs"
	"golang.org/x/tools/go/packages"
)

// FindFlows finds FlowDev flows in the given packages and returns the
// functions or methods containing them.
func FindFlows(pkgs []*packages.Package) ([]*ast.FuncDecl, error) {
	var err error
	flows := make([]*ast.FuncDecl, 0, 4096)

	for _, pkg := range pkgs {
		flows, err = addFlowsFromPackage(flows, pkg)
		if err != nil {
			return nil, err
		}
	}
	return flows, nil
}
func addFlowsFromPackage(flows []*ast.FuncDecl, pkg *packages.Package,
) ([]*ast.FuncDecl, error) {

	if pkgs.IsTestPackage(pkg) {
		return nil, nil
	}

	for _, astf := range pkg.Syntax {
		var err error
		flows, err = addFlowsFromFile(flows, astf)
		if err != nil {
			return nil, err
		}
	}
	return flows, nil
}
func addFlowsFromFile(flows []*ast.FuncDecl, file *ast.File,
) ([]*ast.FuncDecl, error) {
	return flows, nil
}
