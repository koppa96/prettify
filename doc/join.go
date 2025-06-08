package doc

import "io"

// Join is a document element that renders it child nodes with the separating node between them.
type Join struct {
	// The node that separates each child
	Sep Node

	// The child nodes
	Nodes []Node
}

func (j Join) Render(ctx *RenderContext, w io.Writer) error {
	for i, node := range j.Nodes {
		if i > 0 {
			err := j.Sep.Render(ctx, w)
			if err != nil {
				return err
			}
		}

		err := node.Render(ctx, w)
		if err != nil {
			return err
		}
	}

	return nil
}
