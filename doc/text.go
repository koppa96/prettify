package doc

import (
	"fmt"
	"io"
	"unicode/utf8"
)

// TextNode is a document element that renders it string content.
type TextNode struct {
	Value string
	cache lengthCache
}

func (t *TextNode) FlatLength() (int, bool) {
	return t.cache.flatLength(func() (int, bool) {
		return utf8.RuneCountInString(t.Value), true
	})
}

func (t *TextNode) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte(t.Value))
	if err != nil {
		return err
	}

	length, _ := t.FlatLength()
	ctx.CurrentColumn += length
	return nil
}

func Text(value string) *TextNode {
	return &TextNode{
		Value: value,
	}
}

func Textf(format string, args ...any) *TextNode {
	return Text(fmt.Sprintf(format, args...))
}
