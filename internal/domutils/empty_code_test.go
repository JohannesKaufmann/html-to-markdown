package domutils

import (
	"context"
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
)

func TestRemoveEmptyCode(t *testing.T) {
	runs := []struct {
		desc  string
		input string

		expectedBefore string
		expectedAfter  string
	}{
		{
			desc:  "",
			input: `<p>before<code><pre>middle</pre></code>after</p>`,
			expectedBefore: `
├─body
│ ├─p
│ │ ├─#text "before"
│ │ ├─code
│ ├─pre
│ │ ├─code
│ │ │ ├─#text "middle"
│ ├─#text "after"
│ ├─p
			`,
			expectedAfter: `
├─body
│ ├─p
│ │ ├─#text "before"
│ ├─pre
│ │ ├─code
│ │ │ ├─#text "middle"
│ ├─#text "after"
│ ├─p
			`,
		},
		{
			desc:  "two empty code nodes",
			input: `<p><code></code></p>between<p><code></code></p>`,
			expectedBefore: `
├─body
│ ├─p
│ │ ├─code
│ ├─#text "between"
│ ├─p
│ │ ├─code
			`,
			expectedAfter: `
├─body
│ ├─p
│ ├─#text "between"
│ ├─p
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			tester.ExpectRepresentation(t, doc, "before", run.expectedBefore)

			RemoveEmptyCode(context.TODO(), doc)

			tester.ExpectRepresentation(t, doc, "output", run.expectedAfter)
		})
	}

}
