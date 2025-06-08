package doc

import "io"

type Join struct {
	Sep   Node
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
