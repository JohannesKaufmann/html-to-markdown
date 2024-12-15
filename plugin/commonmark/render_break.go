package commonmark

import (
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/marker"
	"golang.org/x/net/html"
)

func (c *commonmark) renderBreak(_ converter.Context, w converter.Writer, _ *html.Node) converter.RenderStatus {
	w.Write(marker.BytesMarkerLineBreak)
	w.Write(marker.BytesMarkerLineBreak)
	return converter.RenderSuccess
}
