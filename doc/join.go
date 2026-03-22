package doc

import "io"

// JoinNode is a document element that renders its child nodes with the separating node between them.
type JoinNode struct {
	// The node that separates each child
	Sep Node

	// The child nodes
	Nodes []Node

	cache lengthCache
}

func (j *JoinNode) FlatLength() (int, bool) {
	return j.cache.flatLength(func() (int, bool) {
		if len(j.Nodes) == 0 {
			return 0, true
		}

		if len(j.Nodes) == 1 {
			return j.Nodes[0].FlatLength()
		}

		sepLength, ok := j.Sep.FlatLength()
		if !ok {
			return 0, false
		}

		length := (len(j.Nodes) - 1) * sepLength
		for _, node := range j.Nodes {
			nodeLength, ok := node.FlatLength()
			if !ok {
				return 0, false
			}

			length += nodeLength
		}

		return length, true
	})
}

func (j *JoinNode) Render(ctx *RenderContext, w io.Writer) error {
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

func Join(nodes []Node, sep Node) *JoinNode {
	return &JoinNode{
		Nodes: nodes,
		Sep:   sep,
	}
}
