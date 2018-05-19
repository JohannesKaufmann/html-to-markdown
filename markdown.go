package md

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func initCommonmarkRules() {
	AddRules(
		Rule{
			Filter: []string{"li"},
			Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
				next := goquery.NodeName(selec.Next())
				parent := selec.Parent()
				index := selec.Index()

				var prefix string
				if parent.Is("ol") {
					prefix = strconv.Itoa(index+1) + ". "
				} else {
					prefix = "- "
				}

				text := prefix + content + "\n"
				if next == "" {
					text += "\n\n"
				}

				// fmt.Println("index:", selec.Index(), name)

				// text := "\n\n- " + content + "\n\n"
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
				// fmt.Println("index:", selec.Index())
				// fmt.Println("next:", selec.Next().Text())

				text := "\n\n" + content + "\n\n"
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
			fmt.Printf("default_rule for %s:'%s' \n", name, res)

			builder.WriteString(res)
			return
		}

		for i := len(r) - 1; i >= 0; i-- {
			rule := r[i]
			content := SelecToMD(domain, s, opt)
			res := rule(content, s, opt)
			if res != nil {
				// fmt.Printf("'%s' \n", *res)
				builder.WriteString(*res)
				return
			}
			fmt.Println(name, "was nil")
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
