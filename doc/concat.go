package doc

import "io"

// Concat is a document element that renders multiple Nodes after each other.
type Concat []Node

func (c Concat) Render(ctx *RenderContext, w io.Writer) error {
	for _, node := range c {
		err := node.Render(ctx, w)
		if err != nil {
			return err
		}
	}

	return nil
}
