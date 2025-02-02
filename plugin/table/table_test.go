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
		// - - - - - - - - - - default - - - - - - - - - - //
		{
			desc: "default",
			options: []option{
				WithMergeContentReplication(false),
			},
			input: `
<table>
  <tr>
    <td>A</td>
    <td colspan="3">B</td>
  </tr>
</table>
			`,
			expected: `
|   |   |  |  |
|---|---|--|--|
| A | B |  |  |
			`,
		},

		// - - - - - - - - - - colspan - - - - - - - - - - //
		{
			desc: "colspan=3",
			options: []option{
				WithMergeContentReplication(true),
			},
			input: `
<table>
  <tr>
    <td>A</td>
    <td colspan="3">B</td>
  </tr>
</table>
			`,
			expected: `
|   |   |   |   |
|---|---|---|---|
| A | B | B | B |
			`,
		},
		// - - - - - - - - - - rowspan - - - - - - - - - - //
		{
			desc: "rowspan=3",
			options: []option{
				WithMergeContentReplication(true),
			},
			input: `
<table>
	<tr>
		<td>A</td>
		<td rowspan="3">B</td>
	</tr>
</table>
			`,
			expected: `
|   |   |
|---|---|
| A | B |
|   | B |
|   | B |
			`,
		},

		// - - - - - - - - - - colspan & rowspan - - - - - - - - - - //
		{
			desc: "cell with colspan and rowspan",
			options: []option{
				WithMergeContentReplication(true),
			},
			input: `
<table>
	<tr>
		<td>A</td>
		<td colspan="3" rowspan="3">B</td>
		<td>C</td>
	</tr>
</table>
			`,
			expected: `
|   |   |   |   |   |
|---|---|---|---|---|
| A | B | B | B | C |
|   | B | B | B |   |
|   | B | B | B |   |
			`,
		},
		{
			desc: "shifting content",
			options: []option{
				WithMergeContentReplication(true),
			},
			input: `
<table>
	<tr>
		<td>A</td>
		<td colspan="3" rowspan="3">B</td>
		<td>C</td>
	</tr>
	<tr>
		<td>1</td>
		<td>2</td>
		<td>3</td>
	</tr>
</table>
			`,
			expected: `
|   |   |   |   |   |   |
|---|---|---|---|---|---|
| A | B | B | B | C |   |
| 1 | B | B | B | 2 | 3 |
|   | B | B | B |   |   |
			`,
		},
		{
			desc: "rowspans overlap with colspans",
			options: []option{
				WithMergeContentReplication(true),
			},
			input: `
<table>
	<tr>
		<td rowspan="3">A</td>
		<td colspan="2">B</td>
		<td>C</td>
	</tr>
	<tr>
		<td rowspan="2" colspan="2">D</td>
		<td>E</td>
	</tr>
	<tr>
		<td>F</td>
	</tr>
</table>
			`,
			expected: `
|   |   |   |   |
|---|---|---|---|
| A | B | B | C |
| A | D | D | E |
| A | D | D | F |
			`,
		},
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
