package dupblock

import (
	"bytes"
	"io"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/donglei1234/platform/services/common/jsonx"
)

var (
	lineSeparator = []byte{'\n'}
	cmdSeparator  = []byte{' '}
)

type TextReader struct {
	block       []byte
	blockHeader *reflect.SliceHeader

	lines [][]byte
	index int
}

func NewTextReader(opts ...TextOption) (Reader, error) {
	var data []byte
	o := textOptions{}
	if err := o.build(opts...); err != nil {
		return nil, err
	} else if o.buffer != nil {
		data = o.buffer.Bytes()
	} else if o.data != nil {
		data = o.data
	} else {
		return nil, ErrDataRequired
	}
	return &TextReader{
		block:       data,
		blockHeader: (*reflect.SliceHeader)(unsafe.Pointer(&data)),
		lines:       bytes.Split(data, lineSeparator),
	}, nil
}

func (r *TextReader) getBlockOffset(s []byte) int {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))

	if hdr.Data >= r.blockHeader.Data && (hdr.Data < r.blockHeader.Data+uintptr(r.blockHeader.Len)) {
		return int(hdr.Data - r.blockHeader.Data)
	}

	return -1
}

func (r *TextReader) readJson() (json []byte, ok bool) {
	start := r.lines[r.index]
	end := start

	for r.index < len(r.lines) {
		line := r.lines[r.index]
		r.index++

		if `*` == string(line) {
			end = line
			break
		}
	}

	if startIdx, endIdx := r.getBlockOffset(start), r.getBlockOffset(end); endIdx > startIdx && startIdx > -1 {
		return r.block[startIdx : endIdx-1], true
	} else {
		return nil, false
	}
}

func (r *TextReader) Read(command *Command) error {
	for r.index < len(r.lines) {
		line := r.lines[r.index]
		r.index++

		if len(line) == 0 {
			continue
		} else if string(line) == "set ." && string(r.lines[r.index]) == "*" {
			*command = Command{
				Action: ActionSet,
				To:     ".",
			}
			r.index++
			return nil
		}

		cmdLine := bytes.Split(line, cmdSeparator)

		if cmdLength := len(cmdLine); cmdLength > 0 {
			cmd := ActionFromText(string(cmdLine[0]))

			switch cmd {
			case ActionSetKey:
				*command = Command{
					Action: cmd,
					To:     string(cmdLine[1]),
				}
				return nil

			case ActionSet, ActionPushFront, ActionPushBack, ActionAddUnique, ActionInsert:
				if cmdLength == 2 {
					*command = Command{
						Action: cmd,
						To:     string(cmdLine[1]),
					}

					if jsonBlock, ok := r.readJson(); ok {
						if err := jsonx.Unmarshal(jsonBlock, &command.Value); err != nil {
							return err
						} else {
							return nil
						}
					} else {
						return ErrMalformedCommand
					}
				} else {
					return ErrMalformedCommand
				}

			case ActionCopy, ActionMove, ActionSwap:
				if cmdLength == 3 {
					*command = Command{
						Action: cmd,
						From:   string(cmdLine[1]),
						To:     string(cmdLine[2]),
					}
					return nil
				} else {
					return ErrMalformedCommand
				}

			case ActionIncrement:
				if cmdLength == 3 {
					if delta, err := strconv.Atoi(string(cmdLine[2])); err != nil {
						return ErrMalformedCommand
					} else {
						*command = Command{
							Action: cmd,
							To:     string(cmdLine[1]),
							Delta:  delta,
						}
						return nil
					}
				} else {
					return ErrMalformedCommand
				}

			case ActionDelete:
				if cmdLength == 2 {
					*command = Command{
						Action: cmd,
						To:     string(cmdLine[1]),
					}
					return nil
				} else {
					return ErrMalformedCommand
				}

			default:
				return ErrMalformedCommand
			}
		} else {
			return ErrMalformedCommand
		}
	}

	return io.EOF
}

func (r *TextReader) Rewind() {
	r.index = 0
}
