package doc

import (
	"errors"
	"io"
)

var ErrCannotRenderFlat = errors.New("this node cannot render in flat mode")

// Expand encapsulates a node that cannot be rendered flat. It causes all of its parent groups to render expanded.
type Expand struct {
	Node Node
}

func (e Expand) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		return ErrCannotRenderFlat
	}

	return e.Node.Render(ctx, w)
}
