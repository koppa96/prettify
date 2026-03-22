package doc

import (
	"io"
)

// Group is a document element that controls whether
// its child should be rendered flat or expanded.
type Group struct {
	// The child of the group
	Node Node
}

func (g Group) FlatLength() (int, bool) {
	return g.Node.FlatLength()
}

func (g Group) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		return g.Node.Render(ctx, w)
	}

	flatLength, ok := g.Node.FlatLength()
	if !ok {
		return g.Node.Render(ctx, w)
	}

	if ctx.CurrentColumn+flatLength > ctx.Config.PrintWidth {
		return g.Node.Render(ctx, w)
	}

	flatCtx := WithFlat(ctx, true)
	err := g.Node.Render(flatCtx, w)
	if err != nil {
		return err
	}

	ctx.CurrentColumn = flatCtx.CurrentColumn
	return nil
}
