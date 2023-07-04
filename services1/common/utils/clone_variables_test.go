package utils

type TestInterface interface {
	TestFuncOne()
	TestFuncTwo()
	TestFuncThree()
}

type TestStruct struct {
	Name string
	Age  int
}

func (t TestStruct) TestFuncOne()   {}
func (t TestStruct) TestFuncTwo()   {}
func (t TestStruct) TestFuncThree() {}
