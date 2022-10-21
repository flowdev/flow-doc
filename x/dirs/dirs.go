package dirs

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindRoot finds the root of a project.
// It looks at the following things (highest priority first):
// - The given directory (used unless empty)
// - It looks for go.mod via `go env GOMOD`
// - It looks for a `vendor` directory.
func FindRoot(dir string, ignoreVendor bool) (string, error) {
	if dir != "" {
		return dir, nil
	}

	if ignoreVendor {
		return crawlUpAndFindDirOf(".", "go.mod")
	}
	return crawlUpAndFindDirOf(".", "go.mod", "vendor")
}

func crawlUpAndFindDirOf(startDir string, files ...string) (string, error) {
	absDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("unable to find absolute directory (for %q): %w", startDir, err)
	}
	volName := filepath.VolumeName(absDir)
	oldDir := "" // set to impossible value first!

	for ; absDir != volName && absDir != oldDir; absDir = filepath.Dir(absDir) {
		for _, file := range files {
			path := filepath.Join(absDir, file)
			if _, err = os.Stat(path); err == nil {
				return absDir, nil
			}
		}
		oldDir = absDir
	}
	return "", fmt.Errorf("unable to find root directory for: %s", absDir)
}
