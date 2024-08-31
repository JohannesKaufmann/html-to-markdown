package commonmark_test

import (
	"bytes"
	"testing"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func TestGoldenFiles(t *testing.T) {
	goldenFileConvert := func(htmlInput []byte) ([]byte, error) {
		conv := converter.NewConverter(
			converter.WithPlugins(commonmark.NewCommonmarkPlugin()),
		)

		// This makes the testcases easier to read
		conv.Register.TagStrategy("#comment", converter.StrategyHTMLBlock)

		return conv.ConvertReader(bytes.NewReader(htmlInput))
	}
	roundTripConvert := func(html []byte) (markdown []byte, err error) {
		// For the golden files we are keeping #comment as a block
		// but collapse treats it as an inline element (which it is).
		//
		// So this testcase would cause problems.
		// "<div>before    <!-- -->    after</div>"

		md, err := htmltomarkdown.ConvertString(string(html))

		return []byte(md), err
	}

	tester.GoldenFiles(t, goldenFileConvert, roundTripConvert)
}

func TestOptionFunc(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		options  []commonmark.OptionFunc
		expected string
	}{
		{
			desc: "WithEmDelimiter",
			options: []commonmark.OptionFunc{
				commonmark.WithEmDelimiter("_"),
			},
			input:    `<em>italic</em>`,
			expected: `_italic_`,
		},
		{
			desc: "WithStrongDelimiter",
			options: []commonmark.OptionFunc{
				commonmark.WithStrongDelimiter("__"),
			},
			input:    `<b>bold</b>`,
			expected: `__bold__`,
		},

		{
			desc: "WithHorizontalRule(***)",
			options: []commonmark.OptionFunc{
				commonmark.WithHorizontalRule("***"),
			},
			input:    `<hr />`,
			expected: `***`,
		},
		{
			desc: "WithHorizontalRule(******)",
			options: []commonmark.OptionFunc{
				commonmark.WithHorizontalRule("******"),
			},
			input:    `<hr />`,
			expected: `******`,
		},
		{
			desc: "WithHorizontalRule(---)",
			options: []commonmark.OptionFunc{
				commonmark.WithHorizontalRule("---"),
			},
			input:    `<hr />`,
			expected: `---`,
		},
		{
			desc: "WithHorizontalRule(___)",
			options: []commonmark.OptionFunc{
				commonmark.WithHorizontalRule("___"),
			},
			input:    `<hr />`,
			expected: `___`,
		},

		{
			desc: "WithBulletListMarker(+)",
			options: []commonmark.OptionFunc{
				commonmark.WithBulletListMarker("+"),
			},
			input:    `<ul><li>list item</li></ul>`,
			expected: `+ list item`,
		},
		{
			desc: "WithBulletListMarker(*)",
			options: []commonmark.OptionFunc{
				commonmark.WithBulletListMarker("*"),
			},
			input:    `<ul><li>list a</li></ul>  <ul><li>list b</li></ul>`,
			expected: "* list a\n\n<!--THE END-->\n\n* list b",
		},
		{
			desc: "WithBulletListMarker(*)",
			options: []commonmark.OptionFunc{
				commonmark.WithBulletListMarker("*"),
				commonmark.WithListEndComment(false),
			},
			input:    `<ul><li>list a</li></ul>  <ul><li>list b</li></ul>`,
			expected: "* list a\n\n* list b",
		},

		{
			desc: "WithCodeBlockFence",
			options: []commonmark.OptionFunc{
				commonmark.WithCodeBlockFence("~~~"),
			},
			input:    `<pre><code>hello world</code></pre>`,
			expected: "~~~\nhello world\n~~~",
		},

		{
			desc: "WithHeadingStyle(atx)",
			options: []commonmark.OptionFunc{
				commonmark.WithHeadingStyle("atx"),
			},
			input:    `<h1>important<br/>heading</h1>`,
			expected: "# important heading",
		},
		{
			desc: "WithHeadingStyle(setext)",
			options: []commonmark.OptionFunc{
				commonmark.WithHeadingStyle("setext"),
			},
			input:    `<h1>important<br/>heading</h1>`,
			expected: "important\n\\\nheading\n=========",
		},

		// TODO: handle other link styles
		// {
		// 	desc: "WithLinkStyle(LinkInlined)",
		// 	options: []commonmark.OptionFunc{
		// 		commonmark.WithLinkStyle(commonmark.LinkInlined),
		// 	},
		// 	input:    `<a href="/about">link</a>`,
		// 	expected: "[link](/about)",
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			conv := converter.NewConverter(
				converter.WithPlugins(
					commonmark.NewCommonmarkPlugin(
						tC.options...,
					),
				),
			)

			output, err := conv.ConvertString(tC.input)
			if err != nil {
				t.Error(err)
			}

			if output != tC.expected {
				t.Errorf("expected %q but got %q", tC.expected, output)
			}
		})
	}
}

func TestOptionFunc_ValdationError(t *testing.T) {
	testCases := []struct {
		desc          string
		options       []commonmark.OptionFunc
		expectedError string
	}{
		{
			desc: "WithEmDelimiter(__)",
			options: []commonmark.OptionFunc{
				commonmark.WithEmDelimiter("__"),
			},
			expectedError: `invalid value for EmDelimiter:"__" must be exactly 1 character of "*" or "_"`,
		},
		{
			desc: "WithEmDelimiter(**)",
			options: []commonmark.OptionFunc{
				commonmark.WithEmDelimiter("**"),
			},
			expectedError: `invalid value for EmDelimiter:"**" must be exactly 1 character of "*" or "_"`,
		},

		{
			desc: "WithStrongDelimiter(_)",
			options: []commonmark.OptionFunc{
				commonmark.WithStrongDelimiter("_"),
			},
			expectedError: `invalid value for StrongDelimiter:"_" must be exactly 2 characters of "**" or "__"`,
		},
		{
			desc: "WithStrongDelimiter(*)",
			options: []commonmark.OptionFunc{
				commonmark.WithStrongDelimiter("*"),
			},
			expectedError: `invalid value for StrongDelimiter:"*" must be exactly 2 characters of "**" or "__"`,
		},

		{
			desc: "WithHorizontalRule(* *)",
			options: []commonmark.OptionFunc{
				commonmark.WithHorizontalRule("* *"),
			},
			expectedError: `invalid value for HorizontalRule:"* *" must be at least 3 characters of "*", "_" or "-"`,
		},
		{
			desc: "WithHorizontalRule(+++)",
			options: []commonmark.OptionFunc{
				commonmark.WithHorizontalRule("+++"),
			},
			expectedError: `invalid value for HorizontalRule:"+++" must be at least 3 characters of "*", "_" or "-"`,
		},

		{
			desc: "WithBulletListMarker(_)",
			options: []commonmark.OptionFunc{
				commonmark.WithBulletListMarker("_"),
			},
			expectedError: `invalid value for BulletListMarker:"_" must be one of "-", "+" or "*"`,
		},

		{
			desc: "WithCodeBlockFence(~~)",
			options: []commonmark.OptionFunc{
				commonmark.WithCodeBlockFence("~~"),
			},
			expectedError: "invalid value for CodeBlockFence:\"~~\" must be one of \"```\" or \"~~~\"",
		},

		{
			desc: "WithHeadingStyle(ATX)",
			options: []commonmark.OptionFunc{
				commonmark.WithHeadingStyle("ATX"),
			},
			expectedError: `invalid value for HeadingStyle:"ATX" must be one of "atx" or "setext"`,
		},
		{
			desc: "WithHeadingStyle(misspelling settext)",
			options: []commonmark.OptionFunc{
				commonmark.WithHeadingStyle("settext"),
			},
			expectedError: `invalid value for HeadingStyle:"settext" must be one of "atx" or "setext"`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			conv := converter.NewConverter(
				converter.WithPlugins(
					commonmark.NewCommonmarkPlugin(
						tC.options...,
					),
				),
			)

			_, err := conv.ConvertString("<strong>bold text</strong>")
			if err == nil {
				t.Fatal("expected an error but got nil")
			}

			_, isValidateConfigError := err.(*commonmark.ValidateConfigError)
			if !isValidateConfigError {
				// t.Error("the error is not of type ValidateConfigError")
			}

			actual := err.Error()
			if actual != tC.expectedError {
				t.Errorf("expected %q but got %q", tC.expectedError, actual)
			}
		})
	}
}
