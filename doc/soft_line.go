package doc

import "io"

type SoftLine struct{}

func (SoftLine) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		return nil
	}

	return HardLine{}.Render(ctx, w)
}
