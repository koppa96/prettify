package doc

import (
	"bytes"
	"io"

	"github.com/koppa96/prettify/config"
)

// RenderContext contains the rendering state.
type RenderContext struct {
	// The settings the formatter was invoked with.
	Config config.Config

	// The current indentation level.
	IndentLevel int

	// Whether the currently rendering item should render flat.
	Flat bool

	// The current column position in the line.
	CurrentColumn int
}

var tabBuf = bytes.Repeat([]byte{'\t'}, 10)

func (ctx *RenderContext) WriteIndent(w io.Writer) error {
	if ctx.IndentLevel == 0 {
		ctx.CurrentColumn = 0
		return nil
	}

	if len(tabBuf) < ctx.IndentLevel {
		tabBuf = append(tabBuf, bytes.Repeat([]byte{'\t'}, ctx.IndentLevel-len(tabBuf))...)
	}

	_, err := w.Write(tabBuf[:ctx.IndentLevel])
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
