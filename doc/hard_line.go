package doc

import "io"

type HardLine struct{}

func (HardLine) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	return ctx.WriteIndent(w)
}
