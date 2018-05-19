package md

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func IsInlineElement(name string) bool {
	return true
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
var indentR = regexp.MustCompile(`(?m)\n`)

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
				content = strings.TrimSpace(content)
				content = indentR.ReplaceAllString(content, "\n  ")

				fmt.Printf("li -> '%s' \n", content)

				// content = leadingNewlinesR.ReplaceAllString(content, "")
				// content = trailingNewlinesR.ReplaceAllString(content, "\n")
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
				return &text
			},
		},
		Rule{
			Filter: []string{"p"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				parent := goquery.NodeName(selec.Parent())
				if IsInlineElement(parent) {
					content += "\n"
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
			// log.Fatal("rule not found for", name)
			content := SelecToMD(domain, s, opt)
			res := DefaultRule(content, s, opt)
			// fmt.Printf("default_rule for %s:'%s' \n", name, res)
			fmt.Println(name, "\t-> default rule")

			builder.WriteString(res)
			return
		}

		for i := len(r) - 1; i >= 0; i-- {
			rule := r[i]
			content := SelecToMD(domain, s, opt)
			res := rule(content, s, opt)
			if res != nil {
				// fmt.Println(name, "\t-> not nil")
				// fmt.Printf("'%s' \n", *res)
				builder.WriteString(*res)
				return
			}
			// fmt.Println(name, "\t-> nil")
		}

		// fmt.Println(i, s.Text(), val, ok, name)

		// content := SelecToMD(domain, s)
		// if name == "head" {
		// 	return
		// }
		// if name == "html" || name == "body" || name == "ul" {
		// 	SelecToMD(domain, s)
		// 	return
		// }
		// fmt.Println("name:", name)
		// res := Commonmark.Paragraph.Replacement(s.Text(), s, Options{})
		// fmt.Printf("'%s'", res)
	})
	return builder.String()
}
