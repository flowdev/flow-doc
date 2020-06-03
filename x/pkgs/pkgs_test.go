package pkgs_test

import (
	"testing"

	"github.com/flowdev/ea-flow-doc/x/pkgs"
	"golang.org/x/tools/go/packages"
)

func TestIsTestPackage(t *testing.T) {
	specs := []struct {
		name           string
		givenPkgPath   string
		expectedIsTest bool
	}{
		{
			name:           "standard package",
			givenPkgPath:   "x/config",
			expectedIsTest: false,
		}, {
			name:           "normal test package",
			givenPkgPath:   "x/config_test",
			expectedIsTest: true,
		}, {
			name:           "main test package",
			givenPkgPath:   "cmd/my_service/main.test",
			expectedIsTest: true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			pkg := &packages.Package{
				PkgPath: spec.givenPkgPath,
			}
			actualIsTest := pkgs.IsTestPackage(pkg)
			if actualIsTest != spec.expectedIsTest {
				t.Errorf("expected %t, actual %t", spec.expectedIsTest, actualIsTest)
			}
		})
	}
}
