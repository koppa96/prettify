package doc

import "io"

type Line struct{}

func (Line) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		_, err := w.Write([]byte{' '})
		return err
	}

	return HardLine{}.Render(ctx, w)
}
