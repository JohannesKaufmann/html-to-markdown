package converter

import (
	"bytes"
	"strings"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/collapse"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/domutils"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/textutils"
	"github.com/JohannesKaufmann/html-to-markdown/v2/marker"

	"golang.org/x/net/html"
)

func (conv *Converter) registerBase() {
	conv.Register.TagStrategy("#comment", StrategyRemoveNode)
	conv.Register.TagStrategy("head", StrategyRemoveNode)
	conv.Register.TagStrategy("script", StrategyRemoveNode)
	conv.Register.TagStrategy("style", StrategyRemoveNode)
	conv.Register.TagStrategy("link", StrategyRemoveNode)
	conv.Register.TagStrategy("meta", StrategyRemoveNode)

	conv.Register.TagStrategy("iframe", StrategyRemoveNode)
	conv.Register.TagStrategy("noscript", StrategyRemoveNode)

	conv.Register.TagStrategy("input", StrategyRemoveNode)
	conv.Register.TagStrategy("textarea", StrategyRemoveNode)

	// "tr" is not in the `IsBlockNode` list,
	// but we want to treat is as a block anyway.
	conv.Register.TagStrategy("tr", StrategyMarkdownBlock)

	conv.Register.PreRenderer(conv.preRenderRemove, PriorityEarly)

	// Note: The priority is low, so that collapse runs _after_ all the other functions
	conv.Register.PreRenderer(conv.preRenderCollapse, PriorityLate)

	conv.Register.Renderer(conv.handleRender, PriorityStandard)

	conv.Register.TextTransformer(conv.handleTextTransform, PriorityStandard)

	conv.Register.PostRenderer(conv.postRenderTrimContent, PriorityStandard)
	conv.Register.PostRenderer(conv.postRenderUnescapeContent, PriorityStandard+20)
}

func (conv *Converter) preRenderRemove(ctx Context, doc *html.Node) {
	var finder func(node *html.Node)
	finder = func(node *html.Node) {
		name := dom.NodeName(node)

		if val, _ := conv.getTagStrategy(name); val == StrategyRemoveNode {
			dom.RemoveNode(node)
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			// Because we are sometimes removing a node, this causes problems
			// with the for loop. Using `defer` is a cool trick!
			// https://gist.github.com/loopthrough/17da0f416054401fec355d338727c46e
			defer finder(child)
		}
	}
	finder(doc)

	// - - - - - - - //

	// After removing elements (see above) it can happen that we have
	// two #text nodes right next to each other. This would cause problems
	// with the collapse so we merge them together.
	domutils.MergeAdjacentTextNodes(doc)
}

func (conv *Converter) preRenderCollapse(ctx Context, doc *html.Node) {
	collapse.Collapse(doc)
}

func (conv *Converter) handleRender(ctx Context, w Writer, n *html.Node) RenderStatus {
	name := dom.NodeName(n)

	switch name {
	case "#text":
		return conv.renderText(ctx, w, n)
	}

	return RenderTryNext
}

func (conv *Converter) handleTextTransform(ctx Context, content string) string {

	// TODO: reduce conversion between types
	content = string(conv.escapeContent([]byte(content)))

	return content
}

var characterEntityReplacer = strings.NewReplacer(
	// We are not using `html.EscapeString` because we
	// care about fewer characters
	"<", "&lt;",
	">", "&gt;",
	"&", "&amp;",
)

func (conv *Converter) renderText(ctx Context, w Writer, n *html.Node) RenderStatus {
	content := n.Data

	// TODO: similar to UnEscapers also only escape if nessesary.
	//       "<" only if not followed by space
	//       "&" only if character entity
	content = characterEntityReplacer.Replace(content)

	for _, handler := range conv.getTextTransformHandlers() {
		content = handler.Value(ctx, content)
	}

	w.WriteString(content)
	return RenderSuccess
}

func (conv *Converter) postRenderTrimContent(ctx Context, result []byte) []byte {
	// Remove whitespace from the beginning & end
	result = bytes.TrimFunc(result, marker.IsSpace)

	// Remove too many newlines
	result = textutils.TrimConsecutiveNewlines(result)

	return result
}
func (conv *Converter) postRenderUnescapeContent(ctx Context, result []byte) []byte {
	result = conv.unEscapeContent(result)
	return result
}
