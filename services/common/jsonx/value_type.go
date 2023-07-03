package jsonx

type ValueType int

const (
	ValueTypeUndefined ValueType = iota
	ValueTypeNull      ValueType = iota
	ValueTypeTrue      ValueType = iota
	ValueTypeFalse     ValueType = iota
	ValueTypeString    ValueType = iota
	ValueTypeNumber    ValueType = iota
	ValueTypeObject    ValueType = iota
	ValueTypeArray     ValueType = iota
)

func (v ValueType) IsBasic() bool {
	return v < ValueTypeObject
}

func (v ValueType) CanHaveChildren() bool {
	return v == ValueTypeArray || v == ValueTypeObject
}

func (v ValueType) String() string {
	switch v {
	case ValueTypeUndefined:
		return "undefined"
	case ValueTypeNull:
		return "null"
	case ValueTypeTrue:
		return "true"
	case ValueTypeFalse:
		return "false"
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	default:
		panic("unknown ValueType")
	}
}
