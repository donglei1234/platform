package dupblock

type Writer interface {
	WriteKey(key string) error
	WriteSet(path string, value interface{}) error
	WriteInsert(path string, value interface{}) error
	WriteIncrement(path string, delta int) error
	WritePushFront(path string, value interface{}) error
	WritePushBack(path string, value interface{}) error
	WriteAddUnique(path string, value interface{}) error
	WriteDelete(path string) error
	WriteCopy(srcPath string, dstPath string) error
	WriteMove(srcPath string, dstPath string) error
	WriteSwap(srcPath string, dstPath string) error
}

type Reader interface {
	Read(command *Command) error
	Rewind()
}
