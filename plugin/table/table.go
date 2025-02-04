package table

import (
	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

type option func(p *tablePlugin)

// TODO: comment & better name?
func WithMergeContentReplication(replicate bool) option {
	return func(p *tablePlugin) {
		p.mergeContentReplication = replicate
	}
}

// WithSkipEmptyRows configures the table plugin to omit empty rows from the output.
// An empty row is defined as a row where all cells contain no content or only whitespace.
//
// true = empty rows will be skipped
//
// false = empty rows will be retained
func WithSkipEmptyRows(skip bool) option {
	return func(p *tablePlugin) {
		p.skipEmptyRows = skip
	}
}

// WithHeaderPromotion configures the table plugin to promote the first row to a header
// if no explicit header row is present. If set to true, the first row becomes the header.
func WithHeaderPromotion(promote bool) option {
	return func(p *tablePlugin) {
		p.promoteFirstRowToHeader = promote
	}
}

type tablePlugin struct {
	mergeContentReplication bool
	skipEmptyRows           bool
	promoteFirstRowToHeader bool
}

func NewTablePlugin(opts ...option) converter.Plugin {
	plugin := &tablePlugin{}
	for _, opt := range opts {
		opt(plugin)
	}
	return plugin
}

func (s *tablePlugin) Name() string {
	return "table"
}

func (s *tablePlugin) Init(conv *converter.Converter) error {

	// TODO: register other stuff for table
	// conv.Register.PreRenderer(s.handlePreRender, converter.PriorityMedium)
	// conv.Register.EscapedChar('|')
	// conv.Register.UnEscapers(converter.PriorityMedium, s.handleUnEscapers)

	conv.Register.Renderer(s.handleRender, converter.PriorityStandard)

	return nil
}

func (s *tablePlugin) handleRender(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	name := dom.NodeName(n)
	switch name {
	case "table":
		return s.renderTableBody(ctx, w, n)
	}

	return converter.RenderTryNext
}
