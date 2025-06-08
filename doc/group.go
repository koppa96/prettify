package doc

import (
	"bytes"
	"io"
	"unicode/utf8"
)

type Group struct {
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
