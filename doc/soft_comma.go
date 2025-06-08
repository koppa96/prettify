package doc

import "io"

type SoftComma struct{}

func (SoftComma) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		return nil
	}

	_, err := w.Write([]byte{','})
	return err
}
