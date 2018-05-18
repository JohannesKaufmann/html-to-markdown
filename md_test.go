package md

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var dmp = diffmatchpatch.New()

var tests = []struct {
	name string
	in   string
	out  string
}{
	{"p", `<p>Some Text</p>`, "Some Text"},
	{"p & span", `<p>Some <span>Text</span></p>`, `Some Text`},
	{"h1", `<h1>Text</h1>`, `# Text`},
	{"bold", `<strong>Text</strong>`, `**Text**`},
	{"italic", `<em>Text</em>`, `_Text_`},
	{"strikethrough", `<del>This was mistaken text</del>`, `~~This was mistaken text~~`},
	{"bold and italic", `<strong>This text is <em>extremely</em> important</strong>`, `**This text is _extremely_ important**`},
	{"ul & li",
		`<ul>
			<li>Some Text</li>
			<li>Another Text</li>
		</ul>`,
		`- Some Text
- Another Text`},
	{"ul & li without whitespace",
		`<ul>
			<li>Some Text</li><li>Another Text</li>
		</ul>`,
		`- Some Text
- Another Text`},
	{
		"ul & li with trailing text",
		`
		<p>Some Text</p>
		<ul>
			<li>1</li>
			<li>2</li>
		</ul>
		<p>Some other Text</p>`,
		`Some Text

- 1
- 2

Some other Text`,
	},
	{
		"quote",
		`<p>In the words of Abraham Lincoln:</p>
		<blockquote>Pardon my French</blockquote>`,
		`In the words of Abraham Lincoln:
> Pardon my French`,
	},
	{
		"multiline quote with trailing text",
		`<p>In the words of Abraham Lincoln:</p>
		<blockquote>Pardon my French<br>Something Else</blockquote>
<p>Other Text</p>`,
		`In the words of Abraham Lincoln:
> Pardon my French
> Something Else

Other Text`,
	},
	{
		"link",
		`<p>This site was built using <a href="https://pages.github.com/">GitHub Pages</a>.</p>
		<p>The site <a href="https://pages.github.com/">GitHub Pages</a> is great.</p>`,
		`This site was built using [GitHub Pages](https://pages.github.com/).

The site [GitHub Pages](https://pages.github.com/) is great.`,
	},
	{
		"break long lines",
		`Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.`,
		`Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At
vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren,
no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit
amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut
labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam
et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata
sanctus est Lorem ipsum dolor sit amet.`,
	},
	{
		"inline code",
		`
		<p>Inline <code>code</code> has <code>back-ticks around</code> it.</p>
		`,
		"Inline `code` has `back-ticks around` it.",
	},
	{
		"code block",
		`
<pre><code>No language indicated, so no syntax highlighting. 
But let's throw in a &lt;b&gt;tag&lt;/b&gt;.
</code></pre>`,
		"```\nNo language indicated, so no syntax highlighting. \nBut let's throw in a <b>tag</b>.\n```",
	},

	{
		"Turndown Demo",
		`
<h1>Turndown Demo</h1>

<p>This demonstrates <a href="https://github.com/domchristie/turndown">turndown</a> – an HTML to Markdown converter in JavaScript.</p>

<h2>Usage</h2>

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
		`# Turndown Demo

This demonstrates [turndown](https://github.com/domchristie/turndown) – an
HTML to Markdown converter in JavaScript.

## Usage

---

It aims to be [CommonMark](http://commonmark.org/) compliant, and includes
options to style the output. These options include:

- headingStyle (setext or atx)
- horizontalRule (*, -, or _)
- bullet (*, -, or +)
- codeBlockStyle (indented or fenced)
- fence
- emDelimiter (_ or *)
- strongDelimiter (** or __)
- linkStyle (inlined or referenced)
- linkReferenceStyle (full, collapsed, or shortcut)`,
	},
	{
		"about google",
		`<div class="hero-module">
        <h1 class="w-inner-type-display">“Organize the world’s information and make it universally accessible and useful.”</h1>
        <div class="w-inner-type-subheader-one">
          <p>Since the beginning, our goal has been to develop services that significantly improve the lives of as many people as possible.</p>
<p>Not just for some. For everyone.</p>
        </div>
			</div>`,
		`# “Organize the world’s information and make it universally accessible and useful.”

Since the beginning, our goal has been to develop services that significantly
improve the lives of as many people as possible.

Not just for some. For everyone.`,
	},
	{
		"about google 2",
		`<div class="h-c-tile__body--content">
			<h3>An unconventional company shares its vision</h3>
			<p add-ellipsis="0-100:600-100:1024-165">To mark our 2004 initial public offering, founders Larry Page and Sergey Brin penned what they deemed an “owner’s manual” for Google shareholders.</p>
			<ul class="h-c-tile__links tile-content-links">
				<li class="h-c-tile__link h-c-tile__link--text tile-content-link">
					<a href="https://abc.xyz/investor/founders-letters/2004/ipo-letter.html" class="tile-cta" title="Read our founders’ letter" target="_blank" data-g-category="our company" data-g-action="an unconventional company shares its vision" data-g-label="cta-link" data-g-href="https://abc.xyz/investor/founders-letters/2004/ipo-letter.html">
							Read our founders’<span class="no-wrap"> letter   <svg alt="" class="h-c-tile__link h-c-tile__link--circle" width="20px" height="20px" viewBox="0 0 20 20" version="1.1" xmlns="http://www.w3.org/2000/svg">
			<g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
				<circle fill="#000000" cx="9" cy="9" r="9"></circle>
				<polygon id="Shape" fill="#FFFFFF" points="3 9 4.0575 10.0575 8.25 5.8725 8.25 15 9.75 15 9.75 5.8725 13.935 10.065 15 9 9 3"></polygon>
			</g>
		</svg>
	</span>
			</a>
		</li>
	</ul>
</div>`,
		`### An unconventional company shares its vision

To mark our 2004 initial public offering, founders Larry Page and Sergey Brin
penned what they deemed an “owner’s manual” for Google shareholders.

- [Read our founders’ letter](https://abc.xyz/investor/founders-letters/2004/ipo-letter.html)`,
	},
}

func TestFromString(t *testing.T) {
	/*
			input := `
		      @media screen  {
		 .pane-img-d44741b {
		          background-image: url(https://lh3.googleusercontent.com/TwhRoINtFP2fIw5Q-R_azFmtsDewyZWr4qCMmfwzDvsttAIJMnl90wI7sqKNUE17EjAryLv3wnCCQscvO9CDlsxvJUn1nwIR_zka3Q=w895-h500-l80-sg-rj);
		        }
		      }

		        @media screen  and (min-resolution: 192dpi) {
		 .pane-img-d44741b {
		            background-image: url(https://lh3.googleusercontent.com/TwhRoINtFP2fIw5Q-R_azFmtsDewyZWr4qCMmfwzDvsttAIJMnl90wI7sqKNUE17EjAryLv3wnCCQscvO9CDlsxvJUn1nwIR_zka3Q=w1790-h1000-l80-sg-rj);
		          }
		        }`
			stylesheet, err := parser.Parse(input)
			if err != nil {
				panic("Please fill a bug :)")
			}
			for _, rule := range stylesheet.Rules {
				fmt.Println(rule.EmbedsRules(), rule.Kind == css.QualifiedRule)

			}
			data, _ := json.MarshalIndent(stylesheet, "", "  ")
			fmt.Println(string(data))

			fmt.Print(stylesheet.String())

			fmt.Println("- - - - - -")
	*/

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _, err := FromString("", test.in)
			if err != nil {
				t.Error(err)
			}
			if result != test.out {
				t.Errorf("want:'%s'\ngot:'%s'", test.out, result)
				// t.Errorf("Result not as expected:\n%v\n", diff.LineDiff(result, test.out))

				// diffs := dmp.DiffMain(result, test.out, false)
				// fmt.Println("CHANGES:", dmp.DiffPrettyText(diffs))

				t.Fail()
			}
			// if s != tt.out {
			// 	t.Errorf("got %q, want %q", s, tt.out)
			// }
		})
	}

	// t.Error("REACHED END")
}
