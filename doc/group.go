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

	// Specifies how the group behaves.
	// By default, a group will try to render its contents flat, even if its parent group failed to render flat.
	// If this is set to true, the group will always choose to render in the same mode as its parent.
	Dependent bool
}

func (g Group) Render(ctx *RenderContext, w io.Writer) error {
	if ctx.Flat || g.Dependent {
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
