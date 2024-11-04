package strikethrough

import (
	"bytes"
	"unicode"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/domutils"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/escape"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/textutils"
	"golang.org/x/net/html"
)

type option func(p *strikethroughPlugin)

func WithDelimiter(delimiter string) option {
	return func(p *strikethroughPlugin) {
		p.delimiter = delimiter
	}
}

type strikethroughPlugin struct {
	delimiter string
}

// Strikethrough converts `<strike>`, `<s>`, and `<del>` elements
func NewStrikethroughPlugin(opts ...option) converter.Plugin {
	plugin := &strikethroughPlugin{}
	for _, opt := range opts {
		opt(plugin)
	}

	if plugin.delimiter == "" {
		plugin.delimiter = "~~"
	}

	return plugin
}

func (s *strikethroughPlugin) Name() string {
	return "strikethrough"
}
func (s *strikethroughPlugin) Init(conv *converter.Converter) error {
	conv.Register.PreRenderer(s.handlePreRender, converter.PriorityStandard)

	conv.Register.EscapedChar('~')
	conv.Register.UnEscaper(s.handleUnEscapers, converter.PriorityStandard)

	conv.Register.Renderer(s.handleRender, converter.PriorityStandard)

	return nil
}

func (s *strikethroughPlugin) handlePreRender(ctx converter.Context, doc *html.Node) {
	domutils.RemoveRedundant(doc, nameIsBothStrikethough)
	domutils.MergeAdjacent(doc, nameIsStrikethough)
}

func (s *strikethroughPlugin) handleUnEscapers(chars []byte, index int) int {
	if chars[index] != '~' {
		return -1
	}

	next := escape.GetNextAsRune(chars, index)

	nextIsWhitespace := unicode.IsSpace(next) || next == 0
	if nextIsWhitespace {
		// "not followed by Unicode whitespace"
		return -1
	}

	return 1
}

func nameIsStrikethough(node *html.Node) bool {
	name := dom.NodeName(node)

	return name == "del" || name == "s" || name == "strike"
}
func nameIsBothStrikethough(a *html.Node, b *html.Node) bool {
	return nameIsStrikethough(a) && nameIsStrikethough(b)
}

func (s strikethroughPlugin) handleRender(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	if nameIsStrikethough(n) {
		return s.renderStrikethrough(ctx, w, n)
	}

	return converter.RenderTryNext
}
func (s strikethroughPlugin) renderStrikethrough(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	var buf bytes.Buffer
	ctx.RenderChildNodes(ctx, &buf, n)

	content := buf.Bytes()

	// If there is a newline character between the start and end delimiter
	// the delimiters won't be recognized. Either we remove all newline characters
	// OR on _every_ line we put start & end delimiters.
	content = textutils.DelimiterForEveryLine(content, []byte(s.delimiter))

	w.Write(content)

	return converter.RenderSuccess
}
