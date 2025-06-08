package doc

import "io"

// SoftComma is a document element that renders nothing when rendered in a flat context,
// and renders a comma if rendered in an expanded context.
type SoftComma struct{}

func (SoftComma) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		return nil
	}

	_, err := w.Write([]byte{','})
	return err
}
