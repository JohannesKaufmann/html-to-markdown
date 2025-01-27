package table

import (
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/collapse"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"golang.org/x/net/html"
)

func TestSelectRowNodes(t *testing.T) {
	runs := []struct {
		desc  string
		input string

		expected string
	}{
		{
			desc: "invalid table",
			input: `
<table>
	<tbody>
		<tr>there is no data cell tag</tr>
	</tbody>
</table>
			`,

			// Note: "golang.org/x/net/html" automatically cleans up the "table"
			expected: `
├─body
│ ├─#text "there is no data cell tag"
│ ├─table
│ │ ├─tbody
│ │ │ ├─tr (__test_normal_row__="true")
			`,
		},
		{
			desc:  "completely empty table",
			input: `<table></table>`,

			expected: `
├─body
│ ├─table
			`,
		},
		{
			desc:  "completely empty tbody",
			input: `<table><tbody></tbody></table>`,

			expected: `
├─body
│ ├─table
│ │ ├─tbody
			`,
		},
		{
			desc: "basic table",
			input: `
<table>
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
			// Note: "golang.org/x/net/html" automatically adds the "tbody"
			expected: `
├─body
│ ├─table
│ │ ├─tbody
│ │ │ ├─tr (__test_normal_row__="true")
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "A1"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "A2"
│ │ │ ├─tr (__test_normal_row__="true")
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "B1"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "B2"
			`,
		},
		{
			desc: "basic table with th",
			input: `
<table>
  <tr>
    <th>Heading 1</td>
    <th>Heading 2</td>
  </tr>
  <tr>
    <td>A1</td>
    <td>A2</td>
  </tr>
</table>
			`,
			expected: `
├─body
│ ├─table
│ │ ├─tbody
│ │ │ ├─tr (__test_header_row__="true")
│ │ │ │ ├─th
│ │ │ │ │ ├─#text "Heading 1"
│ │ │ │ ├─th
│ │ │ │ │ ├─#text "Heading 2"
│ │ │ ├─tr (__test_normal_row__="true")
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "A1"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "A2"
			`,
		},
		{
			desc: "with caption, thead, tbody, tfoot",
			input: `
<table>
  <caption>
    A description about the table
  </caption>
  <thead>
    <tr>
      <th scope="col">Name</th>
      <th scope="col">City</th>
      <th scope="col">Age</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <th scope="row">Max Mustermann</th>
      <td>Berlin</td>
      <td>20</td>
    </tr>
    <tr>
      <th scope="row">Max Müller</th>
      <td>München</td>
      <td>30</td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <th scope="row" colspan="2">Average age</th>
      <td>25</td>
    </tr>
  </tfoot>
</table>
			`,
			expected: `
├─body
│ ├─table
│ │ ├─caption
│ │ │ ├─#text "A description about the table"
│ │ ├─thead
│ │ │ ├─tr (__test_header_row__="true")
│ │ │ │ ├─th (scope="col")
│ │ │ │ │ ├─#text "Name"
│ │ │ │ ├─th (scope="col")
│ │ │ │ │ ├─#text "City"
│ │ │ │ ├─th (scope="col")
│ │ │ │ │ ├─#text "Age"
│ │ ├─tbody
│ │ │ ├─tr (__test_normal_row__="true")
│ │ │ │ ├─th (scope="row")
│ │ │ │ │ ├─#text "Max Mustermann"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "Berlin"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "20"
│ │ │ ├─tr (__test_normal_row__="true")
│ │ │ │ ├─th (scope="row")
│ │ │ │ │ ├─#text "Max Müller"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "München"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "30"
│ │ ├─tfoot
│ │ │ ├─tr (__test_normal_row__="true")
│ │ │ │ ├─th (scope="row" colspan="2")
│ │ │ │ │ ├─#text "Average age"
│ │ │ │ ├─td
│ │ │ │ │ ├─#text "25"
			`,
		},

		// TODO: nested table
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			// NOTE FOR FUTURE: I discovered that "golang.org/x/net/html" automatically adds the "tbody".
			// => So we probably don't need to do that much work beforehand.
			collapse.Collapse(doc, nil)

			{
				// We can then see if we correctly *identified* all the necessary table components.
				// For that we add an attribute (just for the test).

				headerRow := selectHeaderRowNode(doc)
				if headerRow != nil {
					headerRow.Attr = append(headerRow.Attr, html.Attribute{
						Key: "__test_header_row__",
						Val: "true",
					})
				}
				for _, n := range selectNormalRowNodes(doc, headerRow) {
					n.Attr = append(n.Attr, html.Attribute{
						Key: "__test_normal_row__",
						Val: "true",
					})
				}
			}

			tester.ExpectRepresentation(t, doc, "output", run.expected)
		})
	}
}
