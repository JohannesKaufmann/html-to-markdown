package md

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var update = flag.Bool("update", false, "update .golden files")

func TestFromString(t *testing.T) {
	var tests = []struct {
		name string

		domain string
		html   string
	}{
		{
			name: "p tag",
			html: `<p>Some Text</p>`,
		},
		{
			name: "two p tags",
			html: `
			<div>
				<p>Text</p>
				<p>Some Text</p>
			</div>
			`,
		},
		{
			name: "span in p tag",
			html: "<p>Some <span>Text</span></p>",
		},
		{
			name: "strong in p tag",
			html: "<p>Some <strong>Text</strong></p>",
		},
		{
			name: "strong in p tag with whitespace",
			html: "<p> Some <strong> Text </strong></p>",
		},
		{
			name: "h1",
			html: "<h1>Header</h1>",
		},
		{
			name: "h2",
			html: "<h2>Header</h2>",
		},
		{
			name: "h6",
			html: "<h6>Header</h6>",
		},
		{
			name: "ul",
			html: `
			<ul>
				<li>Some Thing</li>
				<li>Another Thing</li>
			</ul>
			`,
		},
		{
			name: "ol",
			html: `
			<ol>
				<li>First Thing</li>
				<li>Second Thing</li>
			</ol>
			`,
		},
		{
			name: "indent content in li",
			html: `
			<ul>
				<li>
					Indent First Thing
					<p>Second Thing</p>
				</li>
				<li>Third Thing</li>
			</ul>
			`,
		},
		{
			name: "nested list",
			html: `
			<ul>
				<li>foo
					<ul>
						<li>bar
							<ul>
								<li>baz
									<ul>
										<li>boo</li>
									</ul>
								</li>
							</ul>
						</li>
					</ul>
				</li>
			</ul>
			`,
		},
		{
			name: "ul in ol",
			html: `
			<ol>
				<li>
					<p>First Thing</p>
					<ul>
						<li>Some Thing</li>
						<li>Another Thing</li>
					</ul>
				</li>
				<li>Second Thing</li>
			</ol>
			`,
		},
		{
			name: "empty list item",
			html: `
			<ul>
				<li>foo</li>
				<li></li>
				<li>bar</li>
			</ul>
			`,
		},
		{
			name: "list items ending with a space",
			html: `
<ul>
	<li>List items </li>
	<li>Ending with </li>
	<li>A space </li>
</ul>
			`,
		},
		{
			name: "sup element",
			html: `
			<p>One of the most common equations in all of physics is
			<var>E</var>=<var>m</var><var>c</var><sup>2</sup>.<p>
			`,
		},
		{
			name: "sup element in list",
			html: `
			<p>The ordinal number "fifth" can be abbreviated in various languages as follows:</p>
			<ul>
				<li>English: 5<sup>th</sup></li>
				<li>French: 5<sup>ème</sup></li>
			</ul>
			`,
		},
		{
			name: "image",
			html: `<img alt="website favicon" src="http://commonmark.org/help/images/favicon.png" />`,
		},
		{
			name: "link",
			html: `<a href="http://commonmark.org/">Link</a>`,
		},
		{
			name: "escape strong",
			html: `<p>**Not Strong**
			**Still Not
			Strong**</p>`,
		},
		{
			name: "escape italic",
			html: `<p>_Not Italic_</p>`,
		},
		{
			name: "escape ordered list",
			html: `<p>1. Not List 1. Not List
			1. Not List</p>`,
		},
		{
			name: "escape unordered list",
			html: `<p>- Not List</p>`,
		},
		{
			name: "pre tag",
			html: `
			<div>
				<p>Who ate the most donuts this week?</p>
				<pre><code class="language-foo+bar">Jeff  15
Sam   11
Robin  6</code></pre>
			</div>
			`,
		},
		{
			name: "code tag",
			html: `
			<p>When <code>x = 3</code>, that means <code>x + 2 = 5</code></p>
			`,
		},
		{
			name: "hr",
			html: `
			<p>Some Content</p>
			<hr>
			</p>Other Content</p>
			`,
		},
		{
			name: "blockquote",
			html: `
<blockquote>
Some Quote
Next Line
</blockquote>
			`,
		},
		{
			name: "large blockquote",
			html: `
			<blockquote>
				<p>Allowing an unimportant mistake to pass without comment is a wonderful social grace.</p>
				<p>Ideological differences are no excuse for rudeness.</p>
			</blockquote>
			`,
		},

		{
			name: "turndown demo",
			html: `
			<h1>Turndown Demo</h1>

			<p>This demonstrates <a href="https://github.com/domchristie/turndown">turndown</a> – an HTML to Markdown converter in JavaScript.</p>

			<h2>Usage</h2>

			<pre><code class="language-js">var turndownService = new TurndownService()
console.log(
  turndownService.turndown('&lt;h1&gt;Hello world&lt;/h1&gt;')
)</code></pre>

			<hr />

			<p>It aims to be <a href="http://commonmark.org/">CommonMark</a>
			 compliant, and includes options to style the output. These options include:</p>

			<ul>
			  <li>headingStyle (setext or atx)</li>
			  <li>horizontalRule (*, -, or _)</li>
			  <li>bullet (*, -, or +)</li>
			  <li>codeBlockStyle (indented or fenced)</li>
			  <li>fence</li>
			  <li>emDelimiter (_ or *)</li>
			  <li>strongDelimiter (** or __)</li>
			  <li>linkStyle (inlined or referenced)</li>
			  <li>linkReferenceStyle (full, collapsed, or shortcut)</li>
			</ul>
						`,
		},
		{
			name: "keep tag",
			html: `<keep-tag><p>Content</p></keep-tag>`,
		},
		{
			name: "remove tag",
			html: `<remove-tag><p>Content</p></remove-tag>`,
		},
		{
			/*
				When a header (eg. <h3>) contains any new lines in its body, it will split the header contents
				over multiple lines, breaking the header in Markdown (because in Markdown, a header just
				starts with #'s and anything on the next line is not part of the header). Since in HTML
				and Markdown all white space is treated the same, I chose to replace line endings with spaces.
				-> https://github.com/lunny/html2md/pull/6
			*/
			name: "strip newlines from header",
			html: `
<h3>

Header
Containing

Newlines

</h3>
			`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			converter := NewConverter(test.domain, true, nil)
			converter.Keep("keep-tag").Remove("remove-tag")

			markdown, err := converter.ConvertString(test.html)
			if err != nil {
				t.Error(err)
			}
			data := []byte(markdown)

			// output := blackfriday.Run(data)
			// fmt.Println(string(output))

			gp := filepath.Join("testdata", t.Name()+".golden")
			if *update {
				t.Log("update golden file")
				if err := ioutil.WriteFile(gp, data, 0644); err != nil {
					t.Fatalf("failed to update golden file: %s", err)
				}
			}

			g, err := ioutil.ReadFile(gp)
			if err != nil {
				t.Logf("Result:\n'%s'\n", markdown)
				t.Fatalf("failed reading .golden: %s", err)
			}

			if !bytes.Equal([]byte(markdown), g) {
				t.Errorf("written json does not match .golden file \nexpected:\n'%s'\nbut got:\n'%s'", string(g), markdown)
			}
		})
	}
}

func BenchmarkFromString(b *testing.B) {
	converter := NewConverter("www.google.com", true, nil)

	strongRule := Rule{
		Filter: []string{"strong"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			return nil
		},
	}

	var wg sync.WaitGroup
	convert := func(html string) {
		defer wg.Done()
		_, err := converter.ConvertString(html)
		if err != nil {
			b.Error(err)
		}
	}
	add := func() {
		defer wg.Done()
		converter.AddRules(strongRule)
	}

	for n := 0; n < b.N; n++ {
		wg.Add(2)
		go add()
		go convert("<strong>Bold</strong>")
	}

	wg.Wait()
}
