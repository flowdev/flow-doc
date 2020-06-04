package parse

import (
	"go/ast"
	"reflect"

	"github.com/flowdev/ea-flow-doc/x/pkgs"
	"golang.org/x/tools/go/packages"
)

// PackageFuncs contains all marked functions from parsing a package.
type PackageFuncs struct {
	PkgPath string
	Funcs   []*ast.FuncDecl
}

// FindFlowFuncs finds FlowDev flows in the given packages and returns the
// functions or methods containing them.
func FindFlowFuncs(pkgs []*packages.Package) []PackageFuncs {
	return findMarkedFuncs(pkgs, "//flowdev:flow", true, false)
}

// FindFlowTests finds FlowDev tests in the given packages and returns the
// functions or methods containing them.
func FindFlowTests(pkgs []*packages.Package) []PackageFuncs {
	return findMarkedFuncs(pkgs, "//flowdev:test", false, true)
}

func findMarkedFuncs(pkgs []*packages.Package, mark string, searchProd, searchTest bool,
) []PackageFuncs {
	pkgFuncs := make([]PackageFuncs, 0, 4096)

	for _, pkg := range pkgs {
		pkgFunc := markedFuncsFromPackage(pkg, mark, searchProd, searchTest)
		if len(pkgFunc.Funcs) > 0 {
			pkgFuncs = append(pkgFuncs, pkgFunc)
		}
	}

	return pkgFuncs
}

func markedFuncsFromPackage(pkg *packages.Package, mark string, searchProd, searchTest bool) PackageFuncs {
	if pkgs.IsTestPackage(pkg) {
		if !searchTest {
			return PackageFuncs{}
		}
	} else {
		if !searchProd {
			return PackageFuncs{}
		}
	}

	pkgFunc := PackageFuncs{PkgPath: pkg.PkgPath, Funcs: make([]*ast.FuncDecl, 0, 1024)}
	for _, astf := range pkg.Syntax {
		pkgFunc.Funcs = addMarkedFuncsFromFile(pkgFunc.Funcs, astf, mark)
	}
	return pkgFunc
}

func addMarkedFuncsFromFile(funcs []*ast.FuncDecl, astf *ast.File, mark string) []*ast.FuncDecl {
	for _, decl := range astf.Decls {
		funcs = addMarkedFuncFromDecl(funcs, decl, mark)
	}
	return funcs
}

func addMarkedFuncFromDecl(funcs []*ast.FuncDecl, decl ast.Decl, mark string) []*ast.FuncDecl {
	if isNilInterfaceOrPointer(decl) {
		return funcs
	}

	switch d := decl.(type) {
	case *ast.FuncDecl:
		funcs = addMarkedFuncFromFunc(funcs, d, mark)
	default:
		// marked functions can only be functions or methods
	}
	return funcs
}

func addMarkedFuncFromFunc(funcs []*ast.FuncDecl, fun *ast.FuncDecl, mark string) []*ast.FuncDecl {
	if fun.Doc == nil {
		return funcs
	}

	for _, comm := range fun.Doc.List {
		if comm.Text == mark {
			funcs = append(funcs, fun)
			return funcs
		}
	}

	return funcs
}

func isNilInterfaceOrPointer(v interface{}) bool {
	return v == nil ||
		(reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}
