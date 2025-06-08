package doc

import "io"

// SoftLine is a document element that renders nothing when rendered in a flat context,
// and renders a new line followed by the correct amount of indentation based on the context when rendered in an expanded context.
type SoftLine struct{}

func (SoftLine) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		return nil
	}

	return HardLine{}.Render(ctx, w)
}
