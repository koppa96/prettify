package doc

import (
	"bytes"
	"io"
	"unicode/utf8"
)

// Group is a document element that controls whether
// its child should be rendered flat or expanded.
type Group struct {
	// The child of the group
	Node Node
}

func (g Group) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat {
		return g.Node.Render(ctx, w)
	}

	var buf bytes.Buffer
	err := g.Node.Render(WithFlat(ctx, true), &buf)
	if err != nil {
		return err
	}

	if ctx.CurrentColumn+utf8.RuneCount(buf.Bytes()) > ctx.Config.PrintWidth {
		return g.Node.Render(ctx, w)
	}

	_, err = io.Copy(w, &buf)
	return err
}
