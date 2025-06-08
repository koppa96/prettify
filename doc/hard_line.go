package doc

import "io"

// HardLine is a document element that renders a new line
// followed by the correct amount of indentation based on the context.
type HardLine struct{}

func (HardLine) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	return ctx.WriteIndent(w)
}
