package diff

import (
	"fmt"
	"strings"

	"github.com/donglei1234/platform/services/common/dupblock"
)

type Diff struct {
	Left, Right *Document
	Operations  []DiffOp
}

func (d Diff) DebugDump() {
	fmt.Println("diffs:", len(d.Operations))
	for index, diff := range d.Operations {
		fmt.Printf("%04x: %s\n", index, diff)
	}
}

// DebugString returns a string representation for this diff suitable for debug use.
func (d Diff) DebugString() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintln("diffs:", len(d.Operations)))
	for index, diff := range d.Operations {
		b.WriteString(fmt.Sprintf("%04x: %s\n", index, diff))
	}
	return b.String()
}

func NewDiff(left *Document, right *Document, opts ...Option) Diff {
	operations := make([]DiffOp, 0, 32)

	result := Diff{
		Left:  left,
		Right: right,
	}

	o := NewOptions(opts...)
	includeRoot := o.IsIncludePath(".")

	numLeftFragments, numRightFragments := len(left.fragments), len(right.fragments)

	if numLeftFragments > 0 || numRightFragments > 0 {
		if numLeftFragments == 0 && includeRoot {
			// nothing on left, something on right
			operations = append(operations, DiffOp{
				Type:     DiffOpSet,
				To:       ".",
				Fragment: right.fragments[0],
			})
		} else if numRightFragments == 0 && includeRoot {
			// something on left, nothing on right
			operations = append(operations, DiffOp{
				Type: DiffOpRemove,
				From: ".",
			})
		} else if includeRoot && left.fragments[0].ValueType != right.fragments[0].ValueType {
			// entire document changed type
			operations = append(operations, DiffOp{
				Type:     DiffOpSet,
				To:       ".",
				Fragment: right.fragments[0],
			})
		} else {
			// find all of the new fragments introduced on the right
			newFragments := make([]Fragment, 0, numRightFragments)

			for index := 0; index < numRightFragments; {
				fragment := right.fragments[index]
				index++

				if !o.IsIncludePath(fragment.Path) {
					continue
				}

				if !left.ContainsPath(fragment.Path) {
					newFragments = append(newFragments, fragment)

					// advance past all of the children of this new fragment since they're already included!
					for index < numRightFragments && right.fragments[index].Depth > fragment.Depth {
						index++
					}
				}
			}

			// find all of the deleted and changed fragments
			deletedFragments := make([]Fragment, 0, numLeftFragments)
			changedFragments := make([]Fragment, 0, numLeftFragments)

			for index := 0; index < numLeftFragments; {
				leftFragment := left.fragments[index]
				index++

				if !o.IsIncludePath(leftFragment.Path) {
					continue
				}

				if rightFragment, ok := right.FindFragmentByPath(leftFragment.Path); ok {
					if index > 1 && !leftFragment.IsEqual(rightFragment) && leftFragment.ValueType.IsBasic() && rightFragment.ValueType.IsBasic() {
						changedFragments = append(changedFragments, *rightFragment)
					}
				} else {
					deletedFragments = append(deletedFragments, leftFragment)

					// advance past all of the children of this deleted fragment since they're already included!
					for index < numLeftFragments && left.fragments[index].Depth > leftFragment.Depth {
						index++
					}
				}
			}

			// generate diffs for new fragments
			for _, fragment := range newFragments {
				// @TODO: SNICHOLS: disabling as part of https://github.com/89trillion/platform/services/issues/622
				//if fragment.ValueType == ValueTypeArray || fragment.ValueType == ValueTypeObject {
				//	if oldFragment, ok := left.FindFragmentByHash(fragment.Hash); ok {
				//		operations = append(operations, DiffOp{
				//			Type: DiffOpCopy,
				//			From: oldFragment.Path,
				//			To:   fragment.Path,
				//		})
				//
				//		continue
				//	}
				//}

				// fallthrough case is set
				operations = append(operations, DiffOp{
					Type:     DiffOpSet,
					To:       fragment.Path,
					Fragment: fragment,
				})
			}

			// generate diffs for changed fragments
			for _, fragment := range changedFragments {
				operations = append(operations, DiffOp{
					Type:     DiffOpSet,
					To:       fragment.Path,
					Fragment: fragment,
				})
			}

			// generate diffs for deleted fragments
			for i := len(deletedFragments); i > 0; i-- {
				operations = append(operations, DiffOp{
					Type: DiffOpRemove,
					From: deletedFragments[i-1].Path,
				})
			}
		}
	}

	// store the result
	result.Operations = operations

	return result
}

// Note: On error, no rollback possible.
func (d Diff) WriteDUPBlock(w dupblock.Writer) error {
	for _, op := range d.Operations {
		switch op.Type {
		case DiffOpSet:
			if e := w.WriteSet(op.To, op.Fragment.Value); e != nil {
				return e
			}

		case DiffOpRemove:
			if e := w.WriteDelete(op.From); e != nil {
				return e
			}

		case DiffOpCopy:
			if e := w.WriteCopy(op.From, op.To); e != nil {
				return e
			}

		case DiffOpMove:
			if e := w.WriteMove(op.From, op.To); e != nil {
				return e
			}

		case DiffOpInsert:
			if e := w.WriteInsert(op.To, op.Fragment.Value); e != nil {
				return e
			}

		default:
			return ErrUnknownDiff
		}
	}

	return nil
}
