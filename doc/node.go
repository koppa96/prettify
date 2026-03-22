package doc

import "io"

// Node is an element of a Doc.
type Node interface {
	// FlatLength returns the length of the node if it is rendered in a single line.
	// If the node cannot be rendered in a single line, FlatLengt returns false in the second return variable.
	FlatLength() (int, bool)

	// Render renders the contents of the Node based on the context's state into the writer.
	Render(ctx *RenderContext, w io.Writer) error
}
