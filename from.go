// Package md converts html to markdown.
//
//  converter := md.NewConverter("", true, nil)
//
//  html = `<strong>Important</strong>`
//
//  markdown, err := converter.ConvertString(html)
//  if err != nil {
//    log.Fatal(err)
//  }
//  fmt.Println("md ->", markdown)
// Or if you are already using goquery:
//  markdown, err := converter.Convert(selec)
package md

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type simpleRuleFunc func(content string, selec *goquery.Selection, options *Options) *string
type ruleFunc func(content string, selec *goquery.Selection, options *Options) (res AdvancedResult, skip bool)

// Converter is initialized by NewConverter.
type Converter struct {
	m      sync.RWMutex
	rules  map[string][]ruleFunc
	keep   map[string]struct{}
	remove map[string]struct{}

	Before func(selec *goquery.Selection)

	// TODO: REMOVE!!!
	// dom      *goquery.Selection
	// leading  []string
	// trailing []string
	// Plugin -> ReportError, ... (not public)

	domain  string
	options Options
}

// TODO: some STATE -> naming???
// TODO: should Plugin be called on every convert

func validate(val string, possible ...string) error {
	for _, e := range possible {
		if e == val {
			return nil
		}
	}
	return fmt.Errorf("field must be one of %v but got %s", possible, val)
}
func validateOptions(opt Options) error {
	if err := validate(opt.HeadingStyle, "setext", "atx"); err != nil {
		return err
	}
	if strings.Count(opt.HorizontalRule, "*") < 3 &&
		strings.Count(opt.HorizontalRule, "_") < 3 &&
		strings.Count(opt.HorizontalRule, "-") < 3 {
		return errors.New("HorizontalRule must be at least 3 characters of '*', '_' or '-' but got " + opt.HorizontalRule)
	}

	if err := validate(opt.BulletListMarker, "-", "+", "*"); err != nil {
		return err
	}
	if err := validate(opt.CodeBlockStyle, "indented", "fenced"); err != nil {
		return err
	}
	if err := validate(opt.Fence, "```", "~~~"); err != nil {
		return err
	}
	if err := validate(opt.EmDelimiter, "_", "*"); err != nil {
		return err
	}
	if err := validate(opt.StrongDelimiter, "**", "__"); err != nil {
		return err
	}
	if err := validate(opt.LinkStyle, "inlined", "referenced"); err != nil {
		return err
	}
	if err := validate(opt.LinkReferenceStyle, "full", "collapsed", "shortcut"); err != nil {
		return err
	}

	return nil
}

// NewConverter initializes a new converter and holds all the rules.
// - `domain` is used for links and images to convert relative urls ("/image.png") to absolute urls.
// - CommonMark is the default set of rules. Set enableCommonmark to false if you want
//   to customize everything using AddRules and DONT want to fallback to default rules.
func NewConverter(domain string, enableCommonmark bool, options *Options) *Converter {
	c := &Converter{
		domain: domain,
		rules:  make(map[string][]ruleFunc),
		keep:   make(map[string]struct{}),
		remove: make(map[string]struct{}),
	}

	if enableCommonmark {
		c.AddRules(commonmark...)
		c.remove["script"] = struct{}{}
		c.remove["style"] = struct{}{}
		c.remove["textarea"] = struct{}{}
	}

	// TODO: put domain in options?
	if options == nil {
		options = &Options{}
	}
	if options.HeadingStyle == "" {
		options.HeadingStyle = "atx"
	}
	if options.HorizontalRule == "" {
		options.HorizontalRule = "* * *"
	}
	if options.BulletListMarker == "" {
		options.BulletListMarker = "-"
	}
	if options.CodeBlockStyle == "" {
		options.CodeBlockStyle = "indented"
	}
	if options.Fence == "" {
		options.Fence = "```"
	}
	if options.EmDelimiter == "" {
		options.EmDelimiter = "_"
	}
	if options.StrongDelimiter == "" {
		options.StrongDelimiter = "**"
	}
	if options.LinkStyle == "" {
		options.LinkStyle = "inlined"
	}
	if options.LinkReferenceStyle == "" {
		options.LinkReferenceStyle = "full"
	}

	c.options = *options
	err := validateOptions(c.options)
	if err != nil {
		fmt.Println("markdown options is not valid:", err)
	}

	return c
}
func (c *Converter) getRuleFuncs(tag string) []ruleFunc {
	c.m.RLock()
	defer c.m.RUnlock()

	r, ok := c.rules[tag]
	if !ok || len(r) == 0 {
		if _, keep := c.keep[tag]; keep {
			return []ruleFunc{wrap(ruleKeep)}
		}
		if _, remove := c.remove[tag]; remove {
			return nil // TODO:
		}

		return []ruleFunc{wrap(ruleDefault)}
	}

	return r
}

func wrap(simple simpleRuleFunc) ruleFunc {
	return func(content string, selec *goquery.Selection, opt *Options) (AdvancedResult, bool) {
		res := simple(content, selec, opt)
		if res == nil {
			return AdvancedResult{}, true
		}
		return AdvancedResult{Markdown: *res}, false
	}
}

// AddRules adds the rules that are passed in to the converter.
func (c *Converter) AddRules(rules ...Rule) *Converter {
	c.m.Lock()
	defer c.m.Unlock()

	for _, rule := range rules {
		if len(rule.Filter) == 0 {
			panic("you need to specify at least one filter for your rule")
		}
		for _, filter := range rule.Filter {
			r, _ := c.rules[filter]

			if rule.AdvancedReplacement != nil {
				r = append(r, rule.AdvancedReplacement)
			} else {
				r = append(r, wrap(rule.Replacement))
			}
			c.rules[filter] = r
		}
	}

	return c
}

// Keep certain html tags in the generated output.
func (c *Converter) Keep(tags ...string) *Converter {
	c.m.Lock()
	defer c.m.Unlock()

	for _, tag := range tags {
		c.keep[tag] = struct{}{}
	}
	return c
}

// Remove certain html tags from the source.
func (c *Converter) Remove(tags ...string) *Converter {
	c.m.Lock()
	defer c.m.Unlock()
	for _, tag := range tags {
		c.remove[tag] = struct{}{}
	}
	return c
}

// Plugin can be used to extends functionality beyond what
// is offered by commonmark.
type Plugin func(conv *Converter) []Rule

// Use can be used to add additional functionality to the converter. It is
// used when its not sufficient to use only rules for example in Plugins.
func (c *Converter) Use(plugins ...Plugin) *Converter {
	for _, plugin := range plugins {
		rules := plugin(c)
		c.AddRules(rules...) // TODO: for better perfomance only use one lock for all plugins
	}
	return c
}

// TODO: Find
// TODO: ReportError
// TODO: AddLeading

// Timeout for the http client
var Timeout = time.Second * 10
var netClient = &http.Client{
	Timeout: Timeout,
}

// DomainFromURL removes the path from the url.
func DomainFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	u.Path = ""
	return u.String()
}

var multipleNewLinesRegex = regexp.MustCompile(`[\n]{2,}`)

// Convert returns the content from a goquery selection.
// If you have a goquery document just pass in doc.Selection.
func (c *Converter) Convert(selec *goquery.Selection) string {
	c.m.RLock()
	domain := c.domain
	options := c.options
	l := len(c.rules)
	if l == 0 {
		panic("you have added no rules. either enable commonmark or add you own.")
	}
	c.m.RUnlock()

	selec.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("data-index", strconv.Itoa(i+1))
	})

	res := c.selecToMD(domain, selec, &options)
	markdown := res.Markdown

	if res.Header != "" {
		markdown = res.Header + "\n\n" + markdown
	}
	if res.Footer != "" {
		markdown += "\n\n" + res.Footer
	}

	markdown = strings.TrimSpace(markdown)
	markdown = multipleNewLinesRegex.ReplaceAllString(markdown, "\n\n")

	return markdown
}

// ConvertReader returns the content from a reader and returns a buffer.
func (c *Converter) ConvertReader(reader io.Reader) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return buffer, err
	}
	buffer.WriteString(
		c.Convert(doc.Selection),
	)

	return buffer, nil
}

// ConvertResponse returns the content from a html response.
func (c *Converter) ConvertResponse(res *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	return c.Convert(doc.Selection), nil
}

// ConvertString returns the content from a html string. If you
// already have a goquery selection use `Convert`.
func (c *Converter) ConvertString(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	return c.Convert(doc.Selection), nil
}

// ConvertBytes returns the content from a html byte array.
func (c *Converter) ConvertBytes(bytes []byte) ([]byte, error) {
	res, err := c.ConvertString(string(bytes))
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

// ConvertURL returns the content from the page with that url.
func (c *Converter) ConvertURL(url string) (string, error) {
	// not using goquery.NewDocument directly because of the timeout
	resp, err := netClient.Get(url)
	if err != nil {
		return "", nil
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil
	}
	domain := DomainFromURL(url)
	if c.domain != domain {
		return "", errors.New("expected " + c.domain + " as the domain but got " + domain)
	}
	return c.Convert(doc.Selection), nil
}
