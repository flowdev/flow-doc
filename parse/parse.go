package parse

import (
	"errors"

	"golang.org/x/tools/go/packages"
)

// Dir is parsing a directory (package) and optionally
// the whole directory tree starting at dir.
// All Go packages found are parsed.
func Dir(dir string, tree bool) ([]*packages.Package, error) {
	parseCfg := &packages.Config{
		Logf:  nil, // log.Printf (for debug), nil (for release)
		Dir:   dir,
		Tests: true,
		Mode:  packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedDeps,
	}

	if tree {
		dir += "/..."
	}
	pkgs, err := packages.Load(parseCfg, dir)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("unable to parse packages at: " + dir)
	}
	return pkgs, nil
}
