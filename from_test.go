package md

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var update = flag.Bool("update", false, "update .golden files")

func TestFromString(t *testing.T) {
	/*
		type rule struct {
			Before func()

			After func() AdvancedResult
		}

		type X struct {
			f func(*X)
			i int
		}
		x := X{
			f: func(x *X) {
				x.i = 5
			},
		}
		fmt.Println("before", x.i)
		x.f(&x)
		fmt.Println("after", x.i)
	*/

	var tests = []struct {
		name string

		domain  string
		html    string
		options *Options
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
			name: "setext h1",
			html: "<h1>Header</h1>",
			options: &Options{
				HeadingStyle: "setext",
			},
		},
		{
			name: "setext h2",
			html: "<h2>Header</h2>",
			options: &Options{
				HeadingStyle: "setext",
			},
		},
		{
			name: "setext h3",
			html: "<h3>Header</h3>",
			options: &Options{
				HeadingStyle: "setext",
			},
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
			name: "link with title",
			html: `<a href="http://commonmark.org/" title="Some Text">Link</a>`,
		},
		{
			name: "reference link: full",
			html: `
<a href="http://commonmark.org/first">First Link</a>

<a href="http://commonmark.org/second">Second Link</a>
`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "full",
			},
		},
		{
			name: "reference link: collapsed",
			html: `<a href="http://commonmark.org/">Link</a>`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "collapsed",
			},
		},
		{
			name: "reference link: shortcut",
			html: `<a href="http://commonmark.org/">Link</a>`,
			options: &Options{
				LinkStyle:          "referenced",
				LinkReferenceStyle: "shortcut",
			},
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
		{
			name: "text with whitespace",
			html: `
						<div id="sport_single_post-2" class="widget sport_single_post">
			<h1 class="widget-title">Aktuelles</h1>
			
			<!-- featured image -->
			<div class="mosaic-block fade"><a href="http://www.bonnerruderverein.de/wp-content/uploads/2015/09/BRV-abend.jpg" class="mosaic-overlay fancybox" title="BRV-abend"></a><div class="mosaic-backdrop"><div class="corner-date">25 Mai</div><img src="http://www.bonnerruderverein.de/wp-content/uploads/2015/09/BRV-abend.jpg" alt="" /></div></div>
			<!-- title -->
			<h3 class="title"><a href="http://www.bonnerruderverein.de/bonner-nachtlauf/">9. Bonner Nachtlauf - Einschränkungen am Bootshaus</a></h3>

            <!-- excerpt -->
            am Mittwoch, dem 30. Mai 2018 findet am Bonner Rheinufer der 9. ...
            <a href="http://www.bonnerruderverein.de/bonner-nachtlauf/" class="more">More</a>



			</div>

			<hr />
			
		<div>
			<h1 class="widget-title">Aktuelles</h1>
			<h3 class="title"><a href="some_url">Title</a></h3>

						<!-- excerpt -->
						Fusce dapibus, tellus ac cursus commodo, tortor mauris condimentum nibh, ut fermentum massa justo sit amet risus. Vestibulum id ligula porta felis euismod semper.
						<a href="other_url" class="more">More</a>

		</div>
`,
		},
		{
			name: "pre tag without code tag",
			html: `
<div class="code"><pre>// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {</pre></div>
`,
		},
		/*
					{ // TODO: not working yet
						name: "p tag with lots of whitespace",
						html: `
			<p>
				Sometimes a struct field, function, type, or even a whole package becomes


				redundant or unnecessary, but must be kept for compatibility with existing


				programs.


				To signal that an identifier should not be used, add a paragraph to its doc


				comment that begins with "Deprecated:" followed by some information about the


				deprecation.


				There are a few examples <a href="https://golang.org/search?q=Deprecated:" target="_blank">in the standard library</a>.
			</p>
			`,
					},
		*/
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			converter := NewConverter(test.domain, true, test.options)
			converter.Keep("keep-tag").Remove("remove-tag")

			markdown, err := converter.ConvertString(test.html)
			if err != nil {
				t.Error(err)
			}
			data := []byte(markdown)

			// output := blackfriday.Run(data)
			// fmt.Println(string(output))

			gp := filepath.Join("testdata", strings.ReplaceAll(t.Name(), ":", "")+".golden")
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

func TestAddRules_ChangeContent(t *testing.T) {
	expected := "Some other Content"

	var wasCalled bool
	rule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			wasCalled = true

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}
			return String(expected)
		},
	}

	conv := NewConverter("", true, nil)
	conv.AddRules(rule)
	md, err := conv.ConvertString(`<p>Some Content</p>`)
	if err != nil {
		t.Error(err)
	}

	if md != expected {
		t.Errorf("wanted '%s' but got '%s'", expected, md)
	}
	if !wasCalled {
		t.Error("rule was not called")
	}
}

func TestAddRules_Fallback(t *testing.T) {
	// firstExpected := "Some other Content"
	expected := "Totally different Content"

	var firstWasCalled bool
	var secondWasCalled bool
	firstRule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			firstWasCalled = true
			if secondWasCalled {
				t.Error("expected first rule to be called before second rule. second is already called")
			}

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}

			return nil
		},
	}
	secondRule := Rule{
		Filter: []string{"p"},
		Replacement: func(content string, selec *goquery.Selection, opt *Options) *string {
			secondWasCalled = true
			if !firstWasCalled {
				t.Error("expected first rule to be called before second rule. first is not called yet")
			}

			if content != "Some Content" {
				t.Errorf("got wrong `content`: '%s'", content)
			}
			if !selec.Is("p") {
				t.Error("selec is not p")
			}
			return String(expected)
		},
	}

	conv := NewConverter("", true, nil)
	conv.AddRules(secondRule, firstRule)
	md, err := conv.ConvertString(`<p>Some Content</p>`)
	if err != nil {
		t.Error(err)
	}

	if md != expected {
		t.Errorf("wanted '%s' but got '%s'", expected, md)
	}
	if !firstWasCalled {
		t.Error("first rule was not called")
	}
	if !secondWasCalled {
		t.Error("second rule was not called")
	}
}
func TestWholeSite(t *testing.T) {
	var tests = []struct {
		name   string
		domain string

		file string
	}{
		{
			name: "golang.org",

			domain: "golang.org",
		},
		{
			name:   "bonnerruderverein.de",
			domain: "bonnerruderverein.de",
		},
		{
			name:   "blog.golang.org",
			domain: "blog.golang.org",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			converter := NewConverter(test.domain, true, nil)

			htmlData, err := ioutil.ReadFile(
				filepath.Join("testdata", t.Name()+".html"),
			)
			if err != nil {
				t.Error(err)
			}

			markdownData, err := converter.ConvertBytes(htmlData)
			if err != nil {
				t.Error(err)
			}

			// output := blackfriday.Run(data)
			// fmt.Println(string(output))

			gp := filepath.Join("testdata", t.Name()+".md")
			if *update {
				t.Log("update golden file")
				if err := ioutil.WriteFile(gp, markdownData, 0644); err != nil {
					t.Fatalf("failed to update golden file: %s", err)
				}
			}

			g, err := ioutil.ReadFile(gp)
			if err != nil {
				t.Logf("Result:\n'%s'\n", string(markdownData))
				t.Fatalf("failed reading .golden: %s", err)
			}

			if !bytes.Equal(markdownData, g) {
				t.Errorf("written json does not match .golden file \nexpected:\n'%s'\nbut got:\n'%s'", string(g), string(markdownData))
			}
		})
	}
}
