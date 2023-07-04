package utils

import "os"

var stat = os.Stat

// FileExists checks to see if a path is an existing file, returning true if so.
func FileExists(path string) bool {
	if i, err := stat(path); err != nil {
		return false
	} else {
		return !i.IsDir()
	}
}
