package utils

import (
	"errors"
	"os"
	"time"
)

var (
	Status = &stat

	TestErrFileExistsErr = errors.New("TestErrFileExistsErr")

	FileExists_StatFuncIsDir  = func(string) (os.FileInfo, error) { return &TestFileInfo{amDir: true}, nil }
	FileExists_StatFuncNotDir = func(string) (os.FileInfo, error) { return &TestFileInfo{}, nil }
	FileExists_StatFuncErr    = func(string) (os.FileInfo, error) { return &TestFileInfo{}, TestErrFileExistsErr }
)

type TestFileInfo struct{ amDir bool }

func (t *TestFileInfo) Name() string       { panic("implement me") }
func (t *TestFileInfo) Size() int64        { panic("implement me") }
func (t *TestFileInfo) Mode() os.FileMode  { panic("implement me") }
func (t *TestFileInfo) ModTime() time.Time { panic("implement me") }
func (t *TestFileInfo) IsDir() bool        { return t.amDir }
func (t *TestFileInfo) Sys() interface{}   { panic("implement me") }
