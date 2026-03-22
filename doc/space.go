package doc

import "io"

// Space renders a space character.
type Space struct{}

func (Space) FlatLength() (int, bool) {
	return 0, true
}

func (Space) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{' '})
	if err != nil {
		return err
	}

	ctx.CurrentColumn++

	return nil
}
