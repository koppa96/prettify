package doc

import (
	"io"

	"github.com/koppa96/prettify/config"
)

type Doc struct {
	Node Node
}

func (d *Doc) Render(cfg config.Config, w io.Writer) error {
	ctx := &RenderContext{
		Config: cfg,
	}

	err := d.Node.Render(ctx, w)
	if err != nil {
		return err
	}

	return HardLine{}.Render(ctx, w)
}
