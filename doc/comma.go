package doc

import "io"

// Comma renders a comma character.
type Comma struct{}

func (Comma) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{','})
	return err
}
