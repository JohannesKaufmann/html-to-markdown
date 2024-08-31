package commonmark

import (
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

func (c *commonmark) renderDivider(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {

	w.WriteString("\n\n")
	w.WriteString(c.HorizontalRule)
	w.WriteString("\n\n")

	return converter.RenderSuccess
}
