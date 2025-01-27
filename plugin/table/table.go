package table

import (
	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

type tablePlugin struct {
}

func NewTablePlugin() converter.Plugin {

	return &tablePlugin{}
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
