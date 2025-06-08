package doc

import "io"

// DoubleLine is a document element that renders two new lines,
// followed by the correct amount of indentation based on the context.
type DoubleLine struct{}

func (DoubleLine) Render(ctx *RenderContext, w io.Writer) error {
	_, err := w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	return HardLine{}.Render(ctx, w)
}
