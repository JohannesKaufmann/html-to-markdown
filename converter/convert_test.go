package converter_test

import (
	"testing"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

func TestConvertString(t *testing.T) {
	conv := converter.NewConverter()

	preRenderer := func(ctx converter.Context, doc *html.Node) {
		for _, node := range dom.AllNodes(doc) {
			name := dom.NodeName(node)

			if name == "test" {
				node.Attr[0].Val = "other_value"
			}
		}
	}
	renderer := func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
		name := dom.NodeName(n)

		if name == "#text" {
			w.WriteString(n.Data)
			return converter.RenderSuccess
		} else if name == "test" {
			val := dom.GetAttributeOr(n, "key", "")
			w.WriteString(val)
			return converter.RenderSuccess
		}

		return converter.RenderTryNext
	}
	postRenderer := func(ctx converter.Context, content []byte) []byte {
		return content
	}

	conv.Register.PreRenderer(preRenderer, converter.PriorityStandard)

	conv.Register.Renderer(renderer, converter.PriorityStandard)
	conv.Register.PostRenderer(postRenderer, converter.PriorityStandard)

	output, err := conv.ConvertString(`before<test key="initial_value"></test>after`)
	if err != nil {
		t.Error(err)
	}

	expected := "beforeother_valueafter"
	if output != expected {
		t.Errorf("expected %q but got %q", expected, output)
	}
}

func TestConvertString_ErrNoRenderHandlers(t *testing.T) {
	conv := converter.NewConverter()
	_, err := conv.ConvertString("<strong>bold text</strong>")
	if err == nil {
		t.Fatal("expected an error")
	}
	if err.Error() != "no render handlers are registered. did you forget to register the commonmark plugin?" {
		t.Fatal("expected a different error but got", err)
	}

	// - - - - //

	// Now that we registered something we should not receive an error anymore...
	conv.Register.Renderer(func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
		return converter.RenderTryNext
	}, converter.PriorityStandard)

	_, err = conv.ConvertString("<strong>bold text</strong>")
	if err != nil {
		t.Fatal("did not expect an error since we registered a renderer")
	}
}
