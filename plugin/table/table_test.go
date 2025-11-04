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
			),
		)

		return conv.ConvertReader(bytes.NewReader(htmlInput))
	}

	tester.GoldenFiles(t, goldenFileConvert, goldenFileConvert)
}

func TestOptionFunc_Validation(t *testing.T) {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			NewTablePlugin(
				WithSpanCellBehavior("random"),
			),
		),
	)

	expectedMessage := `error while initializing "table" plugin: unknown value "random" for span cell behavior`
	out, err := conv.ConvertString("<strong>test</strong>")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != expectedMessage {
		t.Errorf("expected %q but got %q", expectedMessage, err.Error())
	}
	if out != "" {
		t.Error("expected empty output")
	}
}

func TestOptionFunc_ColRowSpan(t *testing.T) {
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
				WithSpanCellBehavior(SpanBehaviorEmpty),
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
| A | B |   |   |
			`,
		},

		// - - - - - - - - - - colspan - - - - - - - - - - //
		{
			desc: "colspan=3",
			options: []option{
				WithSpanCellBehavior(SpanBehaviorMirror),
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
				WithSpanCellBehavior(SpanBehaviorMirror),
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
				WithSpanCellBehavior(SpanBehaviorMirror),
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
				WithSpanCellBehavior(SpanBehaviorMirror),
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
				WithSpanCellBehavior(SpanBehaviorMirror),
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

func TestOptionFunc_EmptyRows(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		options  []option
		expected string
	}{
		// - - - - - - - - - - default - - - - - - - - - - //
		{
			desc:    "by default keep empty rows",
			options: []option{},
			input: `
<table>
  <tr>
    <td></td>
    <td>B1</td>
  </tr>
  <tr>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>A3</td>
    <td></td>
  </tr>
</table>
			`,
			expected: `
|    |    |
|----|----|
|    | B1 |
|    |    |
| A3 |    |
			`,
		},
		{
			desc: "some rows are empty",
			options: []option{
				WithSkipEmptyRows(true),
			},
			input: `
<table>
  <tr>
    <td></td>
    <td>B1</td>
  </tr>
  <tr>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>A3</td>
    <td></td>
  </tr>
  <tr>
    <td>    </td>
    <td>    </td>
  </tr>
</table>
			`,
			expected: `
|    |    |
|----|----|
|    | B1 |
| A3 |    |
			`,
		},
		{
			desc: "all rows are empty",
			options: []option{
				WithSkipEmptyRows(true),
			},
			input: `
<p>Before</p>

<table>
  <caption>A description</caption>
  <tr>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td></td>
  </tr>
</table>

<p>After</p>
			`,
			expected: `
Before

A description

After
			`,
		},
		{
			desc: "element that is not rendered",
			options: []option{
				WithSkipEmptyRows(true),
			},
			input: `
<p>Before</p>

<table>
  <tr>
    <td>
      <script type="text/javascript" src="/script"></script>
    </td>
  </tr>
</table>

<p>After</p>
			`,
			expected: `
Before

After
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

func TestOptionFunc_PromoteHeader(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		options  []option
		expected string
	}{
		// - - - - - - - - - - default - - - - - - - - - - //
		{
			desc:    "default",
			options: []option{},
			input: `
<table>
  <tr>
    <td>A1</td>
    <td>B1</td>
  </tr>
  <tr>
    <td>A2</td>
    <td>B2</td>
  </tr>
</table>
			`,
			expected: `
|    |    |
|----|----|
| A1 | B1 |
| A2 | B2 |
			`,
		},
		{
			desc: "not needed",
			options: []option{
				WithHeaderPromotion(true),
			},
			input: `
<table>
  <tr>
    <th>Heading</th>
    <th>Heading</th>
  </tr>
  <tr>
    <td>A1</td>
    <td>B1</td>
  </tr>
  <tr>
    <td>A2</td>
    <td>B2</td>
  </tr>
</table>
			`,
			expected: `
| Heading | Heading |
|---------|---------|
| A1      | B1      |
| A2      | B2      |
			`,
		},

		{
			desc: "promote first row",
			options: []option{
				WithHeaderPromotion(true),
			},
			input: `
<table>
  <tr>
    <td>A1</td>
    <td>B1</td>
  </tr>
  <tr>
    <td>A2</td>
    <td>B2</td>
  </tr>
</table>
			`,
			expected: `
| A1 | B1 |
|----|----|
| A2 | B2 |
			`,
		},
		{
			desc: "promote first row (but it is empty)",
			options: []option{
				WithHeaderPromotion(true),
			},
			input: `
<table>
  <tr>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>A1</td>
    <td>B1</td>
  </tr>
  <tr>
    <td>A2</td>
    <td>B2</td>
  </tr>
</table>
			`,
			expected: `
|    |    |
|----|----|
| A1 | B1 |
| A2 | B2 |
			`,
		},
		{
			desc: "deleted empty rows & promoted first row",
			options: []option{
				WithHeaderPromotion(true),
				WithSkipEmptyRows(true),
			},
			input: `
<table>
  <tr>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>A1</td>
    <td>B1</td>
  </tr>
  <tr>
    <td>A2</td>
    <td>B2</td>
  </tr>
</table>
			`,
			expected: `
| A1 | B1 |
|----|----|
| A2 | B2 |
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

func TestOptionFunc_PresentationTable(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		options  []option
		expected string
	}{
		{
			desc:    "default",
			options: []option{},
			input: `
<table role="presentation">
  <tr>
    <td>A1</td>
    <td>A2</td>
  </tr>
  <tr>
    <td>B1</td>
    <td>B2</td>
  </tr>
</table>
			`,
			expected: `
A1 A2 

B1 B2
			`,
		},
		{
			desc: "keep the presentation table",
			options: []option{
				WithPresentationTables(true),
			},
			input: `
<table role="presentation">
  <tr>
    <td>A1</td>
    <td>A2</td>
  </tr>
  <tr>
    <td>B1</td>
    <td>B2</td>
  </tr>
</table>
			`,
			expected: `
|    |    |
|----|----|
| A1 | A2 |
| B1 | B2 |
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

func TestTableWithNewlines(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		options  []option
		expected string
	}{
		{
			desc: "with skip behavior (default)",
			options: []option{
				WithNewlineBehavior(NewlineBehaviorSkip),
			},
			input: `
<table>
	<tr>
		<td>A11<br />A12</td>
	</tr>
</table>
			`,
			expected: `
A11  
A12
			`,
		},
		{
			desc: "with preserve behavior",
			options: []option{
				WithNewlineBehavior(NewlineBehaviorPreserve),
			},
			input: `
<table>
  <tr>
    <td>A11<br>A12</td>
    <td>B11<br />B12</td>
  </tr>
</table>
			`,
			expected: `
|                |                |
|----------------|----------------|
| A11  <br />A12 | B11  <br />B12 |
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

func TestOptionFunc_PadColumns(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		options  []option
		expected string
	}{
		{
			desc:    "with padding behavior (default)",
			options: []option{},
			input: `
<table>
  <tr>
    <td>This line has some way longer text than the other line below it.</td>
    <td>A2</td>
  </tr>
  <tr>
    <td>B1</td>
    <td>This one has longer text than the line above.</td>
  </tr>
</table>
			`,
			expected: `
|                                                                  |                                               |
|------------------------------------------------------------------|-----------------------------------------------|
| This line has some way longer text than the other line below it. | A2                                            |
| B1                                                               | This one has longer text than the line above. |
			`,
		},
		{
			desc: "without padding behavior",
			options: []option{
				WithCellPadding(CellPaddingNone),
			},
			input: `
<table>
  <tr>
    <td>This line has some way longer text than the other line below it.</td>
    <td>A2</td>
  </tr>
  <tr>
    <td>B1</td>
    <td>This one has longer text than the line above.</td>
  </tr>
</table>
			`,
			expected: `
|||
|---|---|
|This line has some way longer text than the other line below it.|A2|
|B1|This one has longer text than the line above.|
`,
		},
		{
			desc: "with minimal padding behavior",
			options: []option{
				WithCellPadding(CellPaddingMinimal),
			},
			input: `
<table>
  <tr>
    <td>This line has some way longer text than the other line below it.</td>
    <td>A2</td>
  </tr>
  <tr>
    <td>B1</td>
    <td>This one has longer text than the line above.</td>
  </tr>
</table>
			`,
			expected: `
|  |  |
|---|---|
| This line has some way longer text than the other line below it. | A2 |
| B1 | This one has longer text than the line above. |
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
