package doc

import "io"

type Indent struct {
	Node Node
}

func (i Indent) Render(ctx *RenderContext, w io.Writer) error {
	return i.Node.Render(WithIndent(ctx, ctx.IndentLevel+1), w)
}
