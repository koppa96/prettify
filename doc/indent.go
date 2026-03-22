package doc

import "io"

// Indent is a document element that renders its content with increased indentation level.
type Indent struct {
	Node Node
}

func (i Indent) FlatLength() (int, bool) {
	return i.Node.FlatLength()
}

func (i Indent) Render(ctx *RenderContext, w io.Writer) error {
	indentCtx := WithIndent(ctx, ctx.IndentLevel+1)
	err := i.Node.Render(indentCtx, w)
	if err != nil {
		return err
	}

	ctx.CurrentColumn = indentCtx.CurrentColumn
	return nil
}
