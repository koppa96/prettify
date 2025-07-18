package doc

import (
	"fmt"
	"io"
)

// Text is a document element that renders it string content.
type Text string

func (t Text) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte(t))
	if err != nil {
		return err
	}

	ctx.CurrentColumn += len(t)
	return nil
}

func Textf(format string, args ...any) Text {
	return Text(fmt.Sprintf(format, args...))
}
