package converter

import (
	"github.com/JohannesKaufmann/dom"
	"golang.org/x/net/html"
)

func (conv *Converter) handleRenderNode(ctx Context, w Writer, node *html.Node) RenderStatus {
	for _, handler := range conv.getRenderHandlers() {
		status := handler.Value(ctx, w, node)
		if status == RenderSuccess {
			return status
		}
	}

	return conv.fallbackRender(ctx, w, node)
}
func (conv *Converter) handleRenderNodes(ctx Context, w Writer, nodes ...*html.Node) {
	for _, node := range nodes {
		conv.handleRenderNode(ctx, w, node)
	}
}

func (conv *Converter) fallbackRender(ctx Context, w Writer, node *html.Node) RenderStatus {

	name := dom.NodeName(node)
	decision := conv.getTagStrategyWithFallback(name)

	if decision == StrategyHTMLBlockWithMarkdown {
		w.WriteString("<")
		w.WriteString(name)
		// TODO: also render the attributes?
		w.WriteString(">\n\n")

		conv.handleRenderNodes(ctx, w, dom.AllChildNodes(node)...)

		w.WriteString("\n\n</")
		w.WriteString(name)
		w.WriteString(">")
		return RenderSuccess
	}

	if decision == StrategyHTMLBlock {
		w.WriteRune('\n')
		w.WriteRune('\n')
		_ = html.Render(w, node) // TODO: what to do with error?
		w.WriteRune('\n')
		w.WriteRune('\n')
		return RenderSuccess
	}

	if decision == StrategyMarkdownBlock {
		w.WriteRune('\n')
		w.WriteRune('\n')
		conv.handleRenderNodes(ctx, w, dom.AllChildNodes(node)...)
		w.WriteRune('\n')
		w.WriteRune('\n')

		return RenderSuccess
	} else {
		conv.handleRenderNodes(ctx, w, dom.AllChildNodes(node)...)

		return RenderSuccess
	}
}
