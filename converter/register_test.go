package converter_test

import (
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func TestTagType_Priority(t *testing.T) {
	t.Run("the TagType is used in the base RenderAsHTML", func(t *testing.T) {
		input := `<p>This <strong>bold</strong> and <i>italic</i> text</p>`

		conv := converter.NewConverter(
			converter.WithPlugins(
				base.NewBasePlugin(),
				commonmark.NewCommonmarkPlugin(), // <-- registered a renderer for <strong>
			),
		)

		{
			// - - html block - - //
			conv.Register.RendererFor("strong", converter.TagTypeBlock, base.RenderAsHTML, converter.PriorityStandard-10)
			output, err := conv.ConvertString(input)
			if err != nil {
				t.Fatal(err)
			}
			expected := "This \n\n<strong>bold</strong>\n\n and *italic* text"
			if output != expected {
				t.Errorf("expected %q but got %q", expected, output)
			}
		}
		{
			// - - html inline - - //
			conv.Register.RendererFor("strong", converter.TagTypeInline, base.RenderAsHTML, converter.PriorityStandard-20)
			output, err := conv.ConvertString(input)
			if err != nil {
				t.Fatal(err)
			}
			expected := "This <strong>bold</strong> and *italic* text"
			if output != expected {
				t.Errorf("expected %q but got %q", expected, output)
			}
		}
	})
	t.Run("the TagType can override the RenderHandler", func(t *testing.T) {
		conv := converter.NewConverter(
			converter.WithPlugins(
				base.NewBasePlugin(),
				commonmark.NewCommonmarkPlugin(), // <-- registered a renderer for <strong>
			),
		)
		conv.Register.TagType("strong", converter.TagTypeRemove, converter.PriorityStandard)

		// - - - //
		input := `<p>This <strong>bold</strong> and <i>italic</i> text</p>`
		expected := `This and *italic* text`

		output, err := conv.ConvertString(input)
		if err != nil {
			t.Fatal(err)
		}
		if output != expected {
			t.Errorf("expected %q but got %q", expected, output)
		}
	})
	t.Run("the RenderHandler can override the TagType", func(t *testing.T) {
		input := `
			<h1>heading</h1>

			<style>
			h1 { color: red; }
			</style>
		`

		conv := converter.NewConverter(
			converter.WithPlugins(
				base.NewBasePlugin(), // <-- registered the RemoveTagType for <style>
				commonmark.NewCommonmarkPlugin(),
			),
		)

		// - - - default - - - //
		output, err := conv.ConvertString(input)
		if err != nil {
			t.Fatal(err)
		}
		expected1 := "# heading"
		if output != expected1 {
			t.Errorf("expected %q but got %q", expected1, output)
		}

		// - - - overridden (with higher priority) - - - //
		conv.Register.RendererFor("style", converter.TagTypeBlock, base.RenderAsHTML, converter.PriorityEarly)

		output, err = conv.ConvertString(input)
		if err != nil {
			t.Fatal(err)
		}
		expected2 := "# heading\n\n<style>h1 { color: red; }</style>"
		if output != expected2 {
			t.Errorf("expected %q but got %q", expected2, output)
		}
	})
}
