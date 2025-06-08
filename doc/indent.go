package doc

import "io"

// Indent is a document element that renders its content with increased indentation level.
type Indent struct {
	Node Node
}

func (i Indent) Render(ctx *RenderContext, w io.Writer) error {
	return i.Node.Render(WithIndent(ctx, ctx.IndentLevel+1), w)
}
