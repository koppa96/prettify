package doc

import "io"

// Line is a document element that renders a space character if rendered in a flat context,
// or a new line followed by the correct amount of indentation based on the context, if rendered in an expanded context.
type Line struct{}

func (Line) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		_, err := w.Write([]byte{' '})
		return err
	}

	return HardLine{}.Render(ctx, w)
}
