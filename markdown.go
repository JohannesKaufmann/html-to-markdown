package md

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

type Options struct {
	StrongDelimiter string
}

type Rule struct {
	Filter      []string
	Replacement func(content string, selec *goquery.Selection, options *Options) *string
}

var rules map[string][]func(content string, selec *goquery.Selection, options *Options) *string

func init() {
	rules = make(map[string][]func(content string, selec *goquery.Selection, options *Options) *string)
	initCommonmarkRules()
}

var leadingNewlinesR = regexp.MustCompile(`^\n+`)
var trailingNewlinesR = regexp.MustCompile(`\n+$`)

var newlinesR = regexp.MustCompile(`\n+`)
var tabR = regexp.MustCompile(`\t+`)
var indentR = regexp.MustCompile(`(?m)\n`)

var strongItalicR = strings.NewReplacer(
	`*`, `\*`,
	`_`, `\_`,
)
var orderedListR = regexp.MustCompile(`(?m)^(\d+)\.`)
var unorderedListR = regexp.MustCompile(`(?m)^-\s`)

func EscapeMarkdownCharacters(text string) string {
	text = strongItalicR.Replace(text)
	text = orderedListR.ReplaceAllString(text, `$1\.`)
	text = unorderedListR.ReplaceAllString(text, `\- `)
	return text
}
func initCommonmarkRules() {
	AddRules(
		Rule{
			Filter: []string{"ul", "ol"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				fmt.Printf("ul/ol -> '%s' \n", content)

				parent := selec.Parent()
				if parent.Is("li") {
					content = "\n" + content
				} else {
					content = "\n\n" + content + "\n\n"
				}
				return &content
			},
		},
		Rule{
			Filter: []string{"li"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				parent := selec.Parent()
				index := selec.Index()

				var prefix string
				if parent.Is("ol") {
					prefix = strconv.Itoa(index+1) + ". "
				} else {
					prefix = "- "
				}
				content = leadingNewlinesR.ReplaceAllString(content, "")
				content = trailingNewlinesR.ReplaceAllString(content, "\n")

				// content = strings.TrimSpace(content)
				content = indentR.ReplaceAllString(content, "\n    ")

				text := prefix + content + "\n"
				return &text
			},
		},
		Rule{
			Filter: []string{"#text"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				text := selec.Text()
				if trimmed := strings.TrimSpace(text); trimmed == "" {
					return nil
				}
				text = EscapeMarkdownCharacters(text)

				text = newlinesR.ReplaceAllString(text, "")
				text = tabR.ReplaceAllString(text, " ")
				return &text
			},
		},
		Rule{
			Filter: []string{"p"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				parent := goquery.NodeName(selec.Parent())
				if IsInlineElement(parent) || parent == "li" {
					content = "\n" + content + "\n"
					return &content
				}

				content = "\n\n" + content + "\n\n"
				return &content
			},
		},
		Rule{
			Filter: []string{"h1", "h2", "h3", "h4", "h5", "h6"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				node := goquery.NodeName(selec)
				level, err := strconv.Atoi(node[1:])
				if err != nil {
					panic(err)
				}
				prefix := strings.Repeat("#", level)
				text := "\n\n" + prefix + " " + content + "\n\n"
				return &text
			},
		},
		Rule{
			Filter: []string{"strong", "b"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				trimmed := strings.TrimSpace(content)
				if trimmed == "" {
					return &trimmed
				}
				trimmed = opt.StrongDelimiter + trimmed + opt.StrongDelimiter
				return &trimmed
			},
		},
		Rule{
			Filter: []string{"img"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				alt := selec.AttrOr("alt", "")
				src := selec.AttrOr("src", "")
				text := fmt.Sprintf("![%s](%s)", alt, src)
				return &text
			},
		},
		Rule{
			Filter: []string{"a"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				href := selec.AttrOr("href", "")
				text := fmt.Sprintf("[%s](%s)", content, href)
				return &text
			},
		},
	)
}

func AddRules(newRules ...Rule) {
	for _, newRule := range newRules {
		for _, filter := range newRule.Filter {
			r, _ := rules[filter]
			r = append(r, newRule.Replacement)
			rules[filter] = r
		}
	}
}

func DefaultRule(content string, selec *goquery.Selection, opt *Options) string {
	return content
}

func SelecToMD(domain string, selec *goquery.Selection, opt *Options) string {
	var builder strings.Builder
	// TODO: selec.Contents() Children
	// TODO: Text() or DirectText()
	selec.Contents().Each(func(i int, s *goquery.Selection) {
		name := goquery.NodeName(s)
		r, ok := rules[name]
		if !ok {
			content := SelecToMD(domain, s, opt)
			res := DefaultRule(content, s, opt)
			fmt.Println(name, "\t-> default rule")

			builder.WriteString(res)
			return
		}

		for i := len(r) - 1; i >= 0; i-- {
			rule := r[i]
			content := SelecToMD(domain, s, opt)
			res := rule(content, s, opt)
			if res != nil {
				builder.WriteString(*res)
				return
			}
		}
	})
	return builder.String()
}
