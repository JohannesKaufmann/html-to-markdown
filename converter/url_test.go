package converter

import (
	"testing"
)

func TestDefaultAssembleAbsoluteURL(t *testing.T) {
	runs := []struct {
		desc string

		element Element
		input   string
		domain  string

		expected string
	}{
		{
			desc:   "with whitespaces around",
			input:  "  example.com  \n  ",
			domain: "",

			expected: "example.com",
		},
		{
			desc: "empty fragment",

			element: ElementLink,
			input:   "#",
			domain:  "",

			expected: "#",
		},
		{
			desc: "fragment",

			element: ElementLink,
			input:   "#heading",
			domain:  "",

			expected: "#heading",
		},
		{
			desc: "fragment with space",

			element: ElementLink,
			input:   "#my heading",
			domain:  "",

			expected: "#my%20heading",
		},
		{
			desc: "no domain",

			element: ElementLink,
			input:   "/page.html?key=val#hash",
			domain:  "",

			expected: "/page.html?key=val#hash",
		},
		{
			desc: "with domain",

			element: ElementLink,
			input:   "/page.html?key=val#hash",
			domain:  "test.com",

			expected: "http://test.com/page.html?key=val#hash",
		},
		{
			desc: "data uri",

			element: ElementLink,
			input:   "data:image/gif;base64,R0lGODlhEAAQAMQAAORHHOVSKudfOulrSOp3WOyDZu6QdvCchPGolfO0o/XBs/fNwfjZ0frl3/zy7////wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACH5BAkAABAALAAAAAAQABAAAAVVICSOZGlCQAosJ6mu7fiyZeKqNKToQGDsM8hBADgUXoGAiqhSvp5QAnQKGIgUhwFUYLCVDFCrKUE1lBavAViFIDlTImbKC5Gm2hB0SlBCBMQiB0UjIQA7",
			domain:  "test.com",

			expected: "data:image/gif;base64,R0lGODlhEAAQAMQAAORHHOVSKudfOulrSOp3WOyDZu6QdvCchPGolfO0o/XBs/fNwfjZ0frl3/zy7////wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACH5BAkAABAALAAAAAAQABAAAAVVICSOZGlCQAosJ6mu7fiyZeKqNKToQGDsM8hBADgUXoGAiqhSvp5QAnQKGIgUhwFUYLCVDFCrKUE1lBavAViFIDlTImbKC5Gm2hB0SlBCBMQiB0UjIQA7",
		},
		{
			desc: "data uri (with spaces)",

			element: ElementLink,
			input:   "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 56 56' width='56' height='56' %3E%3C/svg%3E",
			domain:  "test.com",

			expected: "data:image/svg+xml,%3Csvg%20xmlns='http://www.w3.org/2000/svg'%20viewBox='0%200%2056%2056'%20width='56'%20height='56'%20%3E%3C/svg%3E",
		},
		{
			desc: "URI scheme",

			element: ElementLink,
			input:   "slack://open?team=abc",
			domain:  "test.com",

			expected: "slack://open?team=abc",
		},

		{
			desc: "already with http",

			element: ElementLink,
			input:   "http://www.example.com",
			domain:  "test.com",

			expected: "http://www.example.com",
		},
		{
			desc: "already with https",

			element: ElementLink,
			input:   "https://www.example.com",
			domain:  "test.com",

			expected: "https://www.example.com",
		},
		{
			desc:   "query parameters",
			input:  "https://www.example.com?a=1&c=2&b=3&x=&y",
			domain: "test.com",

			// Note: If we were to use Query().Encode() the query parameters
			// would be re-ordered as "?a=1&b=3&c=2".
			// We want to keep the original order!
			expected: "https://www.example.com?a=1&c=2&b=3&x=&y",
		},

		{
			desc:   "invalid url with space",
			input:  "https://Open Demo",
			domain: "",

			expected: "https://Open%20Demo",
		},
		{
			desc:   "invalid url with space and brackets",
			input:  "https://Open [foo](uri) Demo",
			domain: "",

			expected: "https://Open%20%5Bfoo%5D%28uri%29%20Demo",
		},

		{
			desc: "mailto",

			element: ElementLink,
			input:   "mailto:hi@example.com?subject=Mail&cc=someoneelse@example.com",
			domain:  "test.com",

			expected: "mailto:hi@example.com?subject=Mail&cc=someoneelse%40example.com",
		},
		{
			desc: "invalid url with newline in mailto",

			element: ElementLink,
			input:   "mailto:hi@example.com?body=Hello\nJohannes",
			domain:  "test.com",

			expected: "mailto:hi@example.com?body=Hello%0AJohannes",
		},
		{
			desc: "mailto with already encoded space",

			element: ElementLink,
			input:   "mailto:hi@example.com?subject=Hello%20Johannes",
			domain:  "test.com",

			expected: "mailto:hi@example.com?subject=Hello%20Johannes",
		},
		{
			desc: "mailto with raw space",

			element: ElementLink,
			input:   "mailto:hi@example.com?subject=Greetings to Johannes",
			domain:  "test.com",

			expected: "mailto:hi@example.com?subject=Greetings%20to%20Johannes",
		},
		{
			desc: "mailto with german 'ä' character",

			element: ElementLink,
			input:   "mailto:hi@example.com?subject=Sie können gern einen Screenshot anhängen",
			domain:  "test.com",

			// Note: While a space " " is allowed inside then <a> href attribute,
			// in markdown the space would cause the link to not be recognized.
			expected: "mailto:hi@example.com?subject=Sie%20k%C3%B6nnen%20gern%20einen%20Screenshot%20anh%C3%A4ngen",
		},
		{
			desc: "mailto with link",

			element: ElementLink,
			input:   "mailto:hi@example.com?body=Article: www.website.com/page.html",
			domain:  "test.com",

			expected: "mailto:hi@example.com?body=Article%3A%20www.website.com%2Fpage.html",
		},
		{
			desc: "brackets inside link #1",

			element: ElementLink,
			input:   "foo(and(bar)",
			domain:  "",

			expected: "foo%28and%28bar%29",
		},
		{
			desc: "brackets inside link #2",

			element: ElementLink,
			input:   "[foo](uri)",
			domain:  "",

			expected: "%5Bfoo%5D%28uri%29",
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			res := defaultAssembleAbsoluteURL(run.element, run.input, run.domain)
			if res != run.expected {
				t.Errorf("expected '%s' but got '%s'", run.expected, res)
			}
		})
	}
}

func TestParseAndEncode(t *testing.T) {
	runs := []struct {
		desc string

		input string

		expected string
	}{
		{
			desc:     "empty string",
			input:    "",
			expected: "",
		},
		{
			desc:     "one pair",
			input:    "a=1",
			expected: "a=1",
		},
		{
			desc:     "multiple pairs",
			input:    "a=1&b=2&c=3",
			expected: "a=1&b=2&c=3",
		},
		{
			desc:     "keep order of multiple pairs",
			input:    "a=1&c=2&b=3",
			expected: "a=1&c=2&b=3",
		},
		{
			desc:     "encode a space",
			input:    "a=hello world&b=hello",
			expected: "a=hello+world&b=hello",
		},

		{
			desc:     "value with space is encoded with percent",
			input:    "key=%20",
			expected: "key=+",
		},
		{
			desc:     "key with space is encoded with percent",
			input:    "%20=value",
			expected: "+=value",
		},
		{
			desc:     "key with space is encoded with plus",
			input:    "key=+",
			expected: "key=+",
		},
		{
			desc:     "value with space is encoded with plus",
			input:    "+=value",
			expected: "+=value",
		},

		{
			desc: "continue on error at value",
			// The error would be:
			//    invalid URL escape "%"
			input:    "a=1&b=%&c=hello world",
			expected: "a=1&b=%&c=hello+world",
		},
		{
			desc: "continue on error at key",
			// The error would be:
			//    invalid URL escape "%"
			input:    "a=1&%=2&c=hello world",
			expected: "a=1&%=2&c=hello+world",
		},
	}

	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			output := ParseAndEncodeQuery(run.input)
			if output != run.expected {
				t.Errorf("expected '%s' but got '%s'", run.expected, output)
			}
		})
	}
}
