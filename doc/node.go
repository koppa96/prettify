package doc

import "io"

type Node interface {
	Render(ctx *RenderContext, w io.Writer) error
}
