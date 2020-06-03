package parse

import (
	"go/ast"

	"github.com/flowdev/ea-flow-doc/x/pkgs"
	"golang.org/x/tools/go/packages"
)

// FindTests finds FlowDev tests in the given packages and returns the
// functions or methods containing them.
func FindTests(pkgs []*packages.Package) ([]*ast.FuncDecl, error) {
	var err error
	tests := make([]*ast.FuncDecl, 0, 4096)

	for _, pkg := range pkgs {
		tests, err = addTestsFromPackage(tests, pkg)
		if err != nil {
			return nil, err
		}
	}
	return tests, nil
}
func addTestsFromPackage(tests []*ast.FuncDecl, pkg *packages.Package,
) ([]*ast.FuncDecl, error) {

	if !pkgs.IsTestPackage(pkg) {
		return nil, nil
	}

	for _, astf := range pkg.Syntax {
		var err error
		tests, err = addTestsFromFile(tests, astf)
		if err != nil {
			return nil, err
		}
	}
	return tests, nil
}
func addTestsFromFile(tests []*ast.FuncDecl, file *ast.File,
) ([]*ast.FuncDecl, error) {
	return tests, nil
}
