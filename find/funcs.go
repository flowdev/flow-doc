package find

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/flowdev/ea-flow-doc/x/pkgs"
	"github.com/flowdev/ea-flow-doc/x/reflect"
	"golang.org/x/tools/go/packages"
)

// PackageFuncs contains all marked functions from parsing a package.
type PackageFuncs struct {
	Fset      *token.FileSet
	TypesInfo *types.Info
	Funcs     []*ast.FuncDecl
}

// FlowFuncs finds FlowDev flows in the given packages and returns the
// functions or methods containing them.
func FlowFuncs(pkgs []*packages.Package) []PackageFuncs {
	return allMarkedFuncs(pkgs, "//flowdev:flow", true, false)
}

// FlowTests finds FlowDev tests in the given packages and returns the
// functions or methods containing them.
func FlowTests(pkgs []*packages.Package) []PackageFuncs {
	return allMarkedFuncs(pkgs, "//flowdev:test", false, true)
}

func allMarkedFuncs(pkgs []*packages.Package, mark string, searchProd, searchTest bool,
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

	pkgFunc := PackageFuncs{Fset: pkg.Fset, TypesInfo: pkg.TypesInfo, Funcs: make([]*ast.FuncDecl, 0, 1024)}
	for _, astf := range pkg.Syntax {
		pkgFunc.Funcs = addMarkedFuncsFromFile(pkgFunc.Funcs, astf, mark)
	}
	//fmt.Println("TYPEs1:", pkg.TypesInfo)
	return pkgFunc
}

func addMarkedFuncsFromFile(funcs []*ast.FuncDecl, astf *ast.File, mark string) []*ast.FuncDecl {
	for _, decl := range astf.Decls {
		funcs = addMarkedFuncFromDecl(funcs, decl, mark)
	}
	return funcs
}

func addMarkedFuncFromDecl(funcs []*ast.FuncDecl, decl ast.Decl, mark string) []*ast.FuncDecl {
	if reflect.IsNilInterfaceOrPointer(decl) {
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
