package parse

import (
	"go/ast"
	"reflect"

	"github.com/flowdev/ea-flow-doc/x/pkgs"
	"golang.org/x/tools/go/packages"
)

// FindFlowFuncs finds FlowDev flows in the given packages and returns the
// functions or methods containing them.
func FindFlowFuncs(pkgs []*packages.Package) []*ast.FuncDecl {
	return findMarkedFuncs(pkgs, "//flowdev:flow", true, false)
}

// FindFlowTests finds FlowDev flows in the given packages and returns the
// functions or methods containing them.
func FindFlowTests(pkgs []*packages.Package) []*ast.FuncDecl {
	return findMarkedFuncs(pkgs, "//flowdev:test", false, true)
}

func findMarkedFuncs(pkgs []*packages.Package, mark string, searchProd, searchTest bool,
) []*ast.FuncDecl {
	flows := make([]*ast.FuncDecl, 0, 4096)

	for _, pkg := range pkgs {
		flows = addMarkedFuncsFromPackage(flows, pkg, mark, searchProd, searchTest)
	}

	return flows
}

func addMarkedFuncsFromPackage(flows []*ast.FuncDecl, pkg *packages.Package, mark string, searchProd, searchTest bool,
) []*ast.FuncDecl {

	if pkgs.IsTestPackage(pkg) {
		if !searchTest {
			return nil
		}
	} else {
		if !searchProd {
			return nil
		}
	}

	for _, astf := range pkg.Syntax {
		flows = addMarkedFuncsFromFile(flows, astf, mark)
	}
	return flows
}

func addMarkedFuncsFromFile(flows []*ast.FuncDecl, astf *ast.File, mark string) []*ast.FuncDecl {
	for _, decl := range astf.Decls {
		flows = addMarkedFuncFromDecl(flows, decl, mark)
	}
	return flows
}

func addMarkedFuncFromDecl(flows []*ast.FuncDecl, decl ast.Decl, mark string) []*ast.FuncDecl {
	if isNilInterfaceOrPointer(decl) {
		return flows
	}

	switch d := decl.(type) {
	case *ast.FuncDecl:
		flows = addMarkedFuncFromFunc(flows, d, mark)
	default:
		// flows can only be implemented by functions or methods
	}
	return flows
}

func addMarkedFuncFromFunc(flows []*ast.FuncDecl, fun *ast.FuncDecl, mark string) []*ast.FuncDecl {
	if fun.Doc == nil {
		return flows
	}

	for _, comm := range fun.Doc.List {
		if comm.Text == mark {
			flows = append(flows, fun)
			return flows
		}
	}

	return flows
}

func isNilInterfaceOrPointer(v interface{}) bool {
	return v == nil ||
		(reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}
