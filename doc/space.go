package doc

import "io"

// Space renders a space character.
type Space struct{}

func (Space) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{' '})
	return err
}
