package diff

import (
	"fmt"
)

type DiffOpType string

const (
	DiffOpSet    DiffOpType = "set"
	DiffOpRemove DiffOpType = "remove"
	DiffOpCopy   DiffOpType = "copy"
	DiffOpMove   DiffOpType = "move"
	DiffOpSwap   DiffOpType = "swap"
	DiffOpInsert DiffOpType = "insert"
)

type DiffOp struct {
	Type     DiffOpType
	From     string
	To       string
	Count    int
	Offset   int
	Fragment Fragment
}

func (d DiffOp) String() string {
	switch d.Type {
	case DiffOpSet:
		return fmt.Sprintf("%s %s to %s", d.Type, d.To, string(d.Fragment.Value))
	case DiffOpRemove:
		return fmt.Sprintf("%s %s %d %d", d.Type, d.From, d.Offset, d.Count)
	case DiffOpCopy, DiffOpMove:
		return fmt.Sprintf("%s %s to %s", d.Type, d.From, d.To)
	case DiffOpSwap:
		return fmt.Sprintf("%s %s with %s", d.Type, d.From, d.To)
	case DiffOpInsert:
		return fmt.Sprintf("%s %s at %s %d %d", d.Type, string(d.Fragment.Value), d.To, d.Offset, d.Count)
	default:
		return fmt.Sprintf("%s from(%s) to(%s) fragment(%s)", d.Type, d.From, d.To, d.Fragment)
	}
}
