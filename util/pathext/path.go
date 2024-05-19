package pathext

import (
	"os"
	"strings"
)

const (
	pkgSep           = "/"
	goModeIdentifier = "go.mod"
)

// JoinPackages calls strings.Join and returns
func JoinPackages(pkgs ...string) string {
	return strings.Join(pkgs, pkgSep)
}

// MkdirIfNotExist makes directories if the input path is not exists
func MkdirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}
func isLink(name string) (bool, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return false, err
	}

	return fi.Mode()&os.ModeSymlink != 0, nil
}
