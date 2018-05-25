package md

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type RuleFunc func(content string, selec *goquery.Selection, options *Options) *string

type Converter struct {
	m      sync.RWMutex
	rules  map[string][]RuleFunc
	keep   map[string]interface{}
	remove map[string]interface{}

	dom      *goquery.Selection
	leading  []string
	trailing []string
	// Plugin -> ReportError, ... (not public)

	domain  string
	options Options
}

func NewConverter(domain string, enableCommonmark bool, options *Options) *Converter {
	c := &Converter{
		domain: domain,
		rules:  make(map[string][]RuleFunc),
		keep:   make(map[string]interface{}),
		remove: make(map[string]interface{}),
	}

	if enableCommonmark {
		c.AddRules(commonmark...)
	}

	// TODO: put domain in options?
	if options == nil {
		options = &Options{}
	}
	if options.StrongDelimiter == "" {
		options.StrongDelimiter = "**"
	}
	if options.Fence == "" {
		options.Fence = "```"
	}
	if options.HR == "" {
		options.HR = "* * *"
	} else {
		fmt.Println("count *", strings.Count(options.HR, "*"))
		// validateOptions()
	}

	c.options = *options
	// ...

	return c
}
func (c *Converter) getRuleFuncs(tag string) []RuleFunc {
	c.m.RLock()
	defer c.m.RUnlock()

	r, ok := c.rules[tag]
	if !ok || len(r) == 0 {
		if _, keep := c.keep[tag]; keep {
			return []RuleFunc{KeepRule}
		}
		if _, remove := c.remove[tag]; remove {
			return nil // TODO:
		}

		return []RuleFunc{DefaultRule}
	}

	return r
}
func (c *Converter) AddRules(rules ...Rule) *Converter {
	c.m.Lock()
	defer c.m.Unlock()

	for _, rule := range rules {
		for _, filter := range rule.Filter {
			r, _ := c.rules[filter]
			r = append(r, rule.Replacement)
			c.rules[filter] = r
		}
	}

	return c
}
func (c *Converter) Keep(tags ...string) *Converter {
	c.m.Lock()
	defer c.m.Unlock()

	for _, tag := range tags {
		c.keep[tag] = 1 // TODO:
	}
	return c
}
func (c *Converter) Remove(tags ...string) *Converter {
	c.m.Lock()
	defer c.m.Unlock()
	for _, tag := range tags {
		c.remove[tag] = 1 // TODO:
	}
	return c
}

// TODO: naming -> Run, Proccess, ...

type Plugin func(conv *Converter) []Rule

func (c *Converter) Use(plugins ...Plugin) *Converter {
	for _, plugin := range plugins {
		rules := plugin(c)
		c.AddRules(rules...) // TODO: for better perfomance only use one lock for all plugins
	}
	return c
}

func (c *Converter) Find(selector string) *goquery.Selection {
	return c.dom.Find(selector)
}
func (c *Converter) ReportError(err error) *Converter {
	// TODO: or maybe channel???
	return c
}
func (c *Converter) AddLeading(text string) *Converter {
	c.m.Lock()
	defer c.m.Unlock()

	c.leading = append(c.leading, text)
	return c
}

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
	c.m.Lock() // DONT NEED THIS?
	c.dom = selec

	domain := c.domain
	options := c.options
	l := len(c.rules)
	if l == 0 {
		panic("you have added no rules. either enable commonmark or add you own.")
	}
	c.m.Unlock()

	markdown := c.selecToMD(domain, selec, &options)

	markdown = strings.TrimSpace(markdown)
	markdown = multipleNewLinesRegex.ReplaceAllString(markdown, "\n\n")

	c.m.RLock()
	markdown = strings.Join(c.leading, "\n") + "\n" + markdown + "\n" + "trailing"
	c.m.RUnlock()

	return markdown
}

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
