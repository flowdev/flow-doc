package pkgs

import (
	"strings"

	"golang.org/x/tools/go/packages"
)

// IsTestPackage returns true if the given package is a test package and false
// otherwise.
func IsTestPackage(pkg *packages.Package) bool {
	result := strings.HasSuffix(pkg.PkgPath, "_test") ||
		strings.HasSuffix(pkg.PkgPath, ".test") ||
		strings.HasSuffix(pkg.ID, ".test]") ||
		strings.HasSuffix(pkg.ID, ".test")
	return result
}
