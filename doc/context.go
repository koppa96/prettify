package doc

import (
	"io"

	"github.com/koppa96/prettify/config"
)

type RenderContext struct {
	Config        config.Config
	IndentLevel   int
	Flat          bool
	CurrentColumn int
	Writer        io.Writer
}

func (ctx *RenderContext) WriteIndent(w io.Writer) error {
	if ctx.IndentLevel == 0 {
		ctx.CurrentColumn = 0
		return nil
	}

	chars := make([]byte, ctx.IndentLevel)
	for i := range chars {
		chars[i] = '\t'
	}

	_, err := w.Write(chars)
	if err != nil {
		return err
	}

	ctx.CurrentColumn = ctx.IndentLevel * ctx.Config.TabWidth
	return nil
}

func WithIndent(ctx *RenderContext, level int) *RenderContext {
	newCtx := *ctx
	newCtx.IndentLevel = level
	return &newCtx
}

func WithFlat(ctx *RenderContext, flat bool) *RenderContext {
	newCtx := *ctx
	newCtx.Flat = flat
	return &newCtx
}
