package doc

import "io"

type DoubleLine struct{}

func (DoubleLine) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	return HardLine{}.Render(ctx, w)
}
