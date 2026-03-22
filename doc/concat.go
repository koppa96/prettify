package doc

import "io"

// ConcatNode is a document element that renders multiple Nodes after each other.
type ConcatNode struct {
	Nodes []Node
	cache lengthCache
}

func (c *ConcatNode) FlatLength() (int, bool) {
	return c.cache.flatLength(func() (int, bool) {
		var length int
		for _, node := range c.Nodes {
			l, ok := node.FlatLength()
			if !ok {
				return 0, false
			}

			length += l
		}

		return length, true
	})
}

func (c *ConcatNode) Render(ctx *RenderContext, w io.Writer) error {
	for _, node := range c.Nodes {
		err := node.Render(ctx, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func Concat(nodes ...Node) *ConcatNode {
	return &ConcatNode{
		Nodes: nodes,
	}
}
