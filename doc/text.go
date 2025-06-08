package doc

import "io"

type Text string

func (t Text) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte(t))
	if err != nil {
		return err
	}

	ctx.CurrentColumn += len(t)
	return nil
}
