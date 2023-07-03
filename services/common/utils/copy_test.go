package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

func testCopyFactory(t *testing.T) (srcPath TempDir, file *os.File, err error) {
	type TestCase struct{ errExpected bool }
	cases := []TestCase{
		{errExpected: false},
		{errExpected: true},
		{errExpected: true},
	}

	for _, tc := range cases {
		if sourceSetup(t); tc.errExpected == true {
			continue
		} else if srcPath, file, err = sourceSetup(t); err != nil && !os.IsNotExist(err) {
			t.Fatal("Couldn't set up the source and destination directories:", err)
		}
	}
	return
}

func sourceSetup(t *testing.T) (srcPath TempDir, sourceFile *os.File, err error) {
	// Directory to be copied from
	if testDir, err := NewTempDir("test_dir"); err != nil {
		t.Fatal("Unable to create temporary directory:", err)
	} else {

		// Set up the source file
		srcFile, err := ioutil.TempFile(testDir.Path(), "test_")
		if !os.IsNotExist(err) && err != nil {
			t.Fatal("Error encountered setting up the source file:", err)
		}
		sourceFile = srcFile
		srcPath = testDir
	}
	return
}

func TestCopyDir(t *testing.T) {
	// Set up source directory
	if srcDir, _, err := testCopyFactory(t); err != nil {
		t.Fatal("Couldn't set up source and destination directories:", err)
	} else {
		defer srcDir.Cleanup()
		// Directory to be copied to
		if copiedDir, err := NewTempDir("copy_to_here"); err != nil {
			t.Fatal("Unable to create temporary directory:", err)
		} else {
			defer copiedDir.Cleanup()
			// Copy the directory
			if err := CopyDir(srcDir.Path(), copiedDir.Path()); err != nil {
				t.Fatal("Couldn't copy the directory", srcDir, ":", err)
			}
		}
	}
}

func TestCopyFile(t *testing.T) {
	// Set up source directory
	if srcPath, srcFile, err := testCopyFactory(t); err != nil {
		t.Fatal("Couldn't set up source and destination directories:", err)
	} else {
		defer srcPath.Cleanup()
		// Copy the file
		if err := CopyFile(srcFile.Name(), srcPath.Path()+"/copiedFile"); err != nil {
			t.Fatal("Couldn't copy the file", srcFile.Name(), ":", err)
		}
	}
}
