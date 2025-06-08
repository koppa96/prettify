package doc

import "io"

// Node is an element of a Doc.
type Node interface {
	// Render renders the contents of the Node based on the context's state into the writer.
	Render(ctx *RenderContext, w io.Writer) error
}
