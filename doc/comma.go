package doc

import "io"

// Comma renders a comma character.
type Comma struct{}

func (Comma) FlatLength() (int, bool) {
	return 1, true
}

func (Comma) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{','})
	if err != nil {
		return err
	}

	ctx.CurrentColumn++

	return nil
}
