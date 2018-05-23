package md

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/escape"
	"github.com/PuerkitoBio/goquery"
)

var commonmark = []Rule{
	Rule{
		Filter: []string{"ul", "ol"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			// fmt.Printf("ul/ol -> '%s' \n", content)

			parent := selec.Parent()
			if parent.Is("li") && parent.Children().Last().IsSelection(selec) {
				// content = "\n" + content
				// panic("ul&li -> parent is li & something")
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
			// remove leading newlines
			content = leadingNewlinesR.ReplaceAllString(content, "")
			// replace trailing newlines with just a single one
			content = trailingNewlinesR.ReplaceAllString(content, "\n")
			// indent
			content = indentR.ReplaceAllString(content, "\n    ")

			// var r = regexp.MustCompile(`\n+`)
			// content = r.ReplaceAllString(content, "\n")
			// fmt.Printf("LI -> '%s' \n", content)
			// content = strings.TrimSpace(content)

			return String(prefix + content + "\n")
		},
	},
	Rule{
		Filter: []string{"#text"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			text := selec.Text()
			if trimmed := strings.TrimSpace(text); trimmed == "" {
				return nil
			}
			// text = newlinesR.ReplaceAllString(text, "")
			text = tabR.ReplaceAllString(text, " ")

			text = escape.Markdown(text)
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
			content = strings.Replace(content, "\n", " ", -1)
			content = strings.Replace(content, "\r", " ", -1)
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
			src, ok := selec.Attr("src")
			if !ok {
				return String("")
			}

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
	Rule{
		Filter: []string{"code"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			text := "`" + content + "`"
			return &text
		},
	},
	Rule{
		Filter: []string{"pre"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			language := selec.Find("code").AttrOr("class", "")
			language = strings.Replace(language, "language-", "", 1)

			text := "\n\n" + opt.Fence + language + "\n" +
				selec.Find("code").Text() +
				"\n" + opt.Fence + "\n\n"
			return &text
		},
	},
	Rule{
		Filter: []string{"hr"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			text := "\n\n" + opt.HR + "\n\n"
			return &text
		},
	},
	Rule{
		Filter: []string{"blockquote"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			// content = strings.Replace(content, "\n", "\n> ->", -1)
			// fmt.Printf("blockquote: '%s' \n\n", content)

			// var r = regexp.MustCompile(`^\n+|\n+$`)
			// content = r.ReplaceAllString(content, "")

			content = strings.TrimSpace(content)
			content = multipleNewLinesRegex.ReplaceAllString(content, "\n\n")

			var beginningR = regexp.MustCompile(`(?m)^`)
			content = beginningR.ReplaceAllString(content, "> ")

			text := "\n\n" + content + "\n\n"
			return &text
		},
	},
}
