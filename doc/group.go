package doc

import (
	"bytes"
	"errors"
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
		if errors.Is(err, ErrCannotRenderFlat) {
			return g.Node.Render(ctx, w)
		}

		return err
	}

	runeCount := utf8.RuneCount(buf.Bytes())
	if ctx.CurrentColumn+runeCount > ctx.Config.PrintWidth {
		return g.Node.Render(ctx, w)
	}

	_, err = io.Copy(w, &buf)
	ctx.CurrentColumn += runeCount

	return err
}
