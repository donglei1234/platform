package dupblock

import (
	"bytes"
	"strconv"

	"github.com/donglei1234/platform/services/common/jsonx"
)

type TextWriter struct {
	buffer *bytes.Buffer
}

func NewTextWriter(opts ...TextOption) (*TextWriter, error) {
	o := textOptions{}
	if err := o.build(opts...); err != nil {
		return nil, err
	} else if o.buffer == nil {
		o.buffer = bytes.NewBuffer(o.data)
	}
	return &TextWriter{
		buffer: o.buffer,
	}, nil
}

func (w *TextWriter) writeTokenWithPathAndJson(token Action, path string, value interface{}) error {
	b := w.buffer
	b.WriteString(token.Text)
	b.WriteByte(' ')
	b.WriteString(path)
	b.WriteByte('\n')
	if data, err := jsonx.Marshal(value); err != nil {
		return err
	} else {
		b.Write(data)
	}
	b.WriteByte('\n')
	return nil
}

func (w *TextWriter) writeTokenWithSrcAndDst(token Action, src string, dst string) error {
	b := w.buffer
	b.WriteString(token.Text)
	b.WriteByte(' ')
	b.WriteString(src)
	b.WriteByte(' ')
	b.WriteString(dst)
	b.WriteByte('\n')
	return nil
}

func (w *TextWriter) Buffer() *bytes.Buffer {
	return w.buffer
}

func (w *TextWriter) WriteKey(key string) error {
	b := w.buffer
	b.WriteString(ActionSetKey.Text)
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('\n')
	return nil
}

func (w *TextWriter) WriteSet(path string, value interface{}) error {
	return w.writeTokenWithPathAndJson(ActionSet, path, value)
}

func (w *TextWriter) WriteInsert(path string, value interface{}) error {
	return w.writeTokenWithPathAndJson(ActionInsert, path, value)
}

func (w *TextWriter) WriteIncrement(path string, delta int) error {
	b := w.buffer
	b.WriteString(ActionIncrement.Text)
	b.WriteByte(' ')
	b.WriteString(path)
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(delta))
	b.WriteByte('\n')
	return nil
}

func (w *TextWriter) WritePushFront(path string, value interface{}) error {
	return w.writeTokenWithPathAndJson(ActionPushFront, path, value)
}

func (w *TextWriter) WritePushBack(path string, value interface{}) error {
	return w.writeTokenWithPathAndJson(ActionPushBack, path, value)
}

func (w *TextWriter) WriteAddUnique(path string, value interface{}) error {
	return w.writeTokenWithPathAndJson(ActionAddUnique, path, value)
}

func (w *TextWriter) WriteDelete(path string) error {
	b := w.buffer
	b.WriteString(ActionDelete.Text)
	b.WriteByte(' ')
	b.WriteString(path)
	b.WriteByte('\n')
	return nil
}

func (w *TextWriter) WriteCopy(srcPath string, dstPath string) error {
	return w.writeTokenWithSrcAndDst(ActionCopy, srcPath, dstPath)
}

func (w *TextWriter) WriteMove(srcPath string, dstPath string) error {
	return w.writeTokenWithSrcAndDst(ActionMove, srcPath, dstPath)
}

func (w *TextWriter) WriteSwap(srcPath string, dstPath string) error {
	return w.writeTokenWithSrcAndDst(ActionSwap, srcPath, dstPath)
}
