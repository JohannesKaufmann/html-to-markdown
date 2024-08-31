package converter_test

import (
	"testing"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

func TestConvertString_Base(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		// - - - - removing nodes - - - - //
		{
			desc: "automatically removed",
			input: `
<div>
	<span>Start</span>
	<script>To be removed</script>
	<span>End</span>
</div>`,
			expected: "Start End",
		},
		{
			desc: "configured to be removed",
			input: `
<div>
	<span>Start</span>
	<my_remove_node>To be removed</my_remove_node>
	<span>End</span>
</div>`,
			expected: "Start End",
		},

		// - - - - markdown block - - - - //
		{
			desc: "automatically a block node",
			input: `
<div>
	<span>Start</span>
	<article>Article</article>
	<span>End</span>
</div>`,
			expected: "Start\n\nArticle\n\nEnd",
		},
		{
			desc: "configured as block node",
			input: `
<div>
	<span>Start</span>
	<my_markdown_block>Block <b>with</b> markdown</my_markdown_block>
	<span>End</span>
</div>`,
			// TODO: expected: "Start\n\nBlock **with** markdown\n\nEnd",
			//       For this the `Collapse` function needs to accept a custom
			//       `isBlockNode` function that gets info from the tag strategies

			expected: "Start \n\nBlock **with** markdown\n\n End",
		},

		// - - - - markdown leaf - - - - //
		{
			desc: "automatically a leaf node",
			input: `
	<div>
		<span>Start</span>
		<span>Span</span>
		<span>End</span>
	</div>`,
			expected: "Start Span End",
		},
		{
			desc: "default a leaf node",
			input: `
	<div>
		<span>Start</span>
		<random>Random</random>
		<span>End</span>
	</div>`,
			expected: "Start Random End",
		},
		{
			desc: "configured as leaf node",
			input: `
	<div>
		<span>Start</span>
		<my_markdown_leaf>Leaf</my_markdown_leaf>
		<span>End</span>
	</div>`,
			expected: "Start Leaf End",
		},
		{
			desc: "overridden to be not removed",
			input: `
<div>
	<span>Start</span>
	<style>Style</style>
	<span>End</span>
</div>`,
			expected: "Start Style End",
		},

		// - - - - keep as html - - - - //
		{
			desc: "configured as html block node",
			input: `
	<div>
		<span>Start</span>
		<my_html_block>  <h1>Test</h1>  </my_html_block>
		<span>End</span>
	</div>`,
			expected: "Start\n\n<my_html_block><h1>Test</h1></my_html_block>\n\nEnd",
		},
		// - - - - html shell with markdown children - - - - //
		{
			desc: "configured as html block with markdown children",
			input: `
		<div>
			<my_html_shell>
				<p><b>bold</b> text</p>
			</my_html_shell>
		</div>`,
			expected: "<my_html_shell>\n\n**bold** text\n\n</my_html_shell>",
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			conv := converter.NewConverter()

			conv.Register.Renderer(func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
				name := dom.NodeName(n)
				if name == "b" {
					w.WriteString("**")
					ctx.RenderChildNodes(ctx, w, n)
					w.WriteString("**")

					return converter.RenderSuccess
				}

				return converter.RenderTryNext
			}, converter.PriorityStandard)

			conv.Register.TagStrategy("my_remove_node", converter.StrategyRemoveNode)
			conv.Register.TagStrategy("my_markdown_block", converter.StrategyMarkdownBlock)
			conv.Register.TagStrategy("my_markdown_leaf", converter.StrategyMarkdownLeaf)
			conv.Register.TagStrategy("my_html_block", converter.StrategyHTMLBlock)
			conv.Register.TagStrategy("my_html_shell", converter.StrategyHTMLBlockWithMarkdown)

			conv.Register.TagStrategy("style", converter.StrategyMarkdownLeaf)

			out, err := conv.ConvertString(run.input)
			if err != nil {
				t.Error(err)
			}
			if out != run.expected {
				t.Errorf("expected %q but got %q", run.expected, out)
			}
		})
	}
}
