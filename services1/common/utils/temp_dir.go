package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TempDir represents a temporary directory suitable for writing stuff into.
type TempDir struct {
	string
}

// NewTempDir creates a new temporary directory under the current user's account with the given prefix.  Use the
// Cleanup function to delete the directory when done.
func NewTempDir(prefix string) (TempDir, error) {
	if d, err := filepath.Abs(
		filepath.Join(
			os.TempDir(),
			prefix,
			fmt.Sprintf("%08x", time.Now()),
		),
	); err != nil {
		return TempDir{}, err
	} else if err := os.MkdirAll(d, os.ModePerm); err != nil {
		return TempDir{}, err
	} else {
		return TempDir{d}, nil
	}
}

// Path returns the absolute path of this temporary directory.
func (t TempDir) Path() string {
	return t.string
}

// Cleanup deletes the temporary directory and all its contents.
func (t *TempDir) Cleanup() error {
	if t.string != "" {
		if err := os.RemoveAll(t.string); err != nil {
			return err
		} else {
			t.string = ""
			return nil
		}
	}
	return nil
}
