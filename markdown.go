package md

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	ruleDefault = func(content string, selec *goquery.Selection, opt *Options) *string {
		return &content
	}
	ruleKeep = func(content string, selec *goquery.Selection, opt *Options) *string {
		element := selec.Get(0)

		var buf bytes.Buffer
		err := html.Render(&buf, element)
		if err != nil {
			panic(err)
		}

		return String(buf.String())
	}
)

var blockElements = []string{
	"address",
	"article",
	"aside",
	"audio",
	"video",
	"blockquote",
	"canvas",
	"dd",
	"div",
	"dl",
	"fieldset",
	"figcaption",
	"figure",
	"footer",
	"form",
	"h1", "h2", "h3", "h4", "h5", "h6",
	"header",
	"hgroup",
	"hr",
	"noscript",
	"ol", "ul",
	"output",
	"p",
	"pre",
	"section",
	"table", "tfoot",
}
var inlineElements = []string{ // -> https://developer.mozilla.org/de/docs/Web/HTML/Inline_elemente
	"b", "big", "i", "small", "tt",
	"abbr", "acronym", "cite", "code", "dfn", "em", "kbd", "strong", "samp", "var",
	"a", "bdo", "br", "img", "map", "object", "q", "script", "span", "sub", "sup",
	"button", "input", "label", "select", "textarea",
}

func IsInlineElement(e string) bool {
	for _, element := range inlineElements {
		if element == e {
			return true
		}
	}
	return false
}
func IsBlockElement(e string) bool {
	for _, element := range blockElements {
		if element == e {
			return true
		}
	}
	return false
}

// String is a helper function to return a pointer.
func String(text string) *string {
	return &text
}

// Options to customize the output. You can change stuff like
// the character that is used for strong text.
type Options struct {
	// "setext" or "atx"
	// default: "atx"
	HeadingStyle string

	// Any Thematic break
	// default: "* * *"
	HorizontalRule string

	// "-", "+", or "*"
	// default: "-"
	BulletListMarker string

	// "indented" or "fenced"
	// default: "indented"
	CodeBlockStyle string

	// ``` or ~~~
	// default: ```
	Fence string

	// _ or *
	// default: _
	EmDelimiter string

	// ** or __
	// default: **
	StrongDelimiter string

	// inlined or referenced
	// default: inlined
	LinkStyle string

	// full, collapsed, or shortcut
	// default: full
	LinkReferenceStyle string
}

type AdvancedResult struct {
	Header   string
	Markdown string
	Footer   string
}

// Rule to convert certain html tags to markdown.
//  md.Rule{
//    Filter: []string{"del", "s", "strike"},
//    Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
//      // You need to return a pointer to a string (md.String is just a helper function).
//      // If you return nil the next function for that html element
//      // will be picked. For example you could only convert an element
//      // if it has a certain class name and fallback if not.
//      return md.String("~" + content + "~")
//    },
//  }
type Rule struct {
	Filter              []string
	Replacement         func(content string, selec *goquery.Selection, options *Options) *string
	AdvancedReplacement func(content string, selec *goquery.Selection, options *Options) (res AdvancedResult, skip bool)
}

var leadingNewlinesR = regexp.MustCompile(`^\n+`)
var trailingNewlinesR = regexp.MustCompile(`\n+$`)

var newlinesR = regexp.MustCompile(`\n+`)
var tabR = regexp.MustCompile(`\t+`)
var indentR = regexp.MustCompile(`(?m)\n`)

func (c *Converter) selecToMD(domain string, selec *goquery.Selection, opt *Options) AdvancedResult {
	var result AdvancedResult

	var builder strings.Builder
	selec.Contents().Each(func(i int, s *goquery.Selection) {
		name := goquery.NodeName(s)
		rules := c.getRuleFuncs(name)

		for i := len(rules) - 1; i >= 0; i-- {
			rule := rules[i]

			content := c.selecToMD(domain, s, opt)
			if content.Header != "" {
				result.Header += content.Header
			}
			if content.Footer != "" {
				result.Footer += content.Footer
			}

			res, skip := rule(content.Markdown, s, opt)
			if res.Header != "" {
				result.Header += res.Header + "\n"
			}
			if res.Footer != "" {
				result.Footer += res.Footer + "\n"
			}

			if !skip {
				builder.WriteString(res.Markdown)
				return
			}
		}
	})
	result.Markdown = builder.String()
	return result
}
