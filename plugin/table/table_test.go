package table

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func TestGoldenFiles(t *testing.T) {
	goldenFileConvert := func(htmlInput []byte) ([]byte, error) {
		conv := converter.NewConverter(
			converter.WithPlugins(
				base.NewBasePlugin(),
				commonmark.NewCommonmarkPlugin(),
				NewTablePlugin(),
			),
		)

		return conv.ConvertReader(bytes.NewReader(htmlInput))
	}

	tester.GoldenFiles(t, goldenFileConvert, goldenFileConvert)
}

func TestOptionFunc(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		options  []option
		expected string
	}{
		// - - - - - - - - - - Colspan - - - - - - - - - - //
		{
			desc: "colspan=3 && WithMergeContentReplication(false)",
			options: []option{
				WithMergeContentReplication(false),
			},
			input: `
<table>
  <tr>
    <td>A</td>
    <td colspan="3">B</td>
    <td>C</td>
  </tr>
</table>
			`,
			expected: `
|   |   |  |  |   |
|---|---|--|--|---|
| A | B |  |  | C |
			`,
		},
		{
			desc: "colspan=3 && WithMergeContentReplication(true)",
			options: []option{
				WithMergeContentReplication(true),
			},
			input: `
<table>
  <tr>
    <td>A</td>
    <td colspan="3">B</td>
    <td>C</td>
  </tr>
</table>
			`,
			expected: `
|   |   |   |   |   |
|---|---|---|---|---|
| A | B | B | B | C |
			`,
		},
		// - - - - - - - - - - Rospan - - - - - - - - - - //
		// TODO: grow the slice
		// 		{
		// 			desc: "rowspan=3 && WithMergeContentReplication(false)",
		// 			options: []option{
		// 				WithMergeContentReplication(false),
		// 			},
		// 			input: `
		// <table>
		//   <tr>
		//     <td>A</td>
		//     <td rowspan="3">B</td>
		//     <td>C</td>
		//   </tr>
		// </table>
		// 			`,
		// 			expected: ``,
		// 		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			conv := converter.NewConverter(
				converter.WithPlugins(
					base.NewBasePlugin(),
					commonmark.NewCommonmarkPlugin(),
					NewTablePlugin(tC.options...),
				),
			)

			output, err := conv.ConvertString(tC.input)
			if err != nil {
				t.Error(err)
			}

			actual := strings.TrimSpace(output)
			expected := strings.TrimSpace(tC.expected)

			if actual != expected {
				t.Errorf("expected\n%s\nbut got\n%s\n", expected, actual)
			}
		})
	}
}
