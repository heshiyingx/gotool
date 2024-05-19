package pathext

import "os"

// FileExists returns true if the specified file is exists.
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
