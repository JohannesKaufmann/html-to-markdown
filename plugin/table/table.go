package table

import (
	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

type option func(p *tablePlugin)

type spanCellBehavior string

const (
	// SpanBehaviorEmpty renders an empty cell.
	SpanBehaviorEmpty spanCellBehavior = "empty"
	// SpanBehaviorMirror renders the same content as the original cell.
	SpanBehaviorMirror spanCellBehavior = "mirror"
)

// WithSpanCellBehavior configures how cells affected by colspan/rowspan attributes
// should be rendered. When a cell spans multiple columns or rows, the affected cells
// can either be empty or contain the same content as the original cell.
func WithSpanCellBehavior(behavior spanCellBehavior) option {
	return func(p *tablePlugin) {
		p.spanCellBehavior = behavior
	}
}

// WithSkipEmptyRows configures the table plugin to omit empty rows from the output.
// An empty row is defined as a row where all cells contain no content or only whitespace.
// When set to true, empty rows will be omitted from the output. When false (default),
// all rows are preserved.
func WithSkipEmptyRows(skip bool) option {
	return func(p *tablePlugin) {
		p.skipEmptyRows = skip
	}
}

// WithHeaderPromotion configures whether the first row should be treated as a header
// when the table has no explicit header row (e.g. <th> elements). When set to true, the
// first row will be converted to a header row with separator dashes. When false (default),
// all rows are treated as regular content.
func WithHeaderPromotion(promote bool) option {
	return func(p *tablePlugin) {
		p.promoteFirstRowToHeader = promote
	}
}

// WithPresentationTables configures whether tables marked with role="presentation"
// should be converted to markdown. When set to true, presentation tables will be
// converted like regular tables. When false (default), these tables are skipped
// since they typically represent layout rather than semantic content.
func WithPresentationTables(convert bool) option {
	return func(p *tablePlugin) {
		p.convertPresentationTables = convert
	}
}

type tablePlugin struct {
	spanCellBehavior          spanCellBehavior
	skipEmptyRows             bool
	promoteFirstRowToHeader   bool
	convertPresentationTables bool
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

	conv.Register.EscapedChar('|')

	conv.Register.Renderer(s.handleRender, converter.PriorityStandard)

	return nil
}

func (s *tablePlugin) handleRender(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	name := dom.NodeName(n)
	switch name {
	case "table":
		return s.renderTable(ctx, w, n)

	case "tr":
		// Normally, when the "table" gets rendered we do NOT go into this case.
		// But as a fallback we separate the rows through newlines.
		return s.renderFallbackRow(ctx, w, n)

	}

	return converter.RenderTryNext
}

func (s *tablePlugin) renderFallbackRow(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	w.WriteString("\n\n")
	ctx.RenderChildNodes(ctx, w, n)
	w.WriteString("\n\n")
	return converter.RenderSuccess
}
