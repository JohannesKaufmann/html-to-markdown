package md_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type GoldenTest struct {
	Name string

	Options map[string]*md.Options
	Plugins []md.Plugin
}

func RunGoldenTest(t *testing.T, tests []GoldenTest) {
	for _, test := range tests {
		if len(test.Options) == 0 {
			test.Options = map[string]*md.Options{
				"default": nil,
			}
		}

		t.Run(test.Name, func(t *testing.T) {
			g := goldie.New(t)

			for key, options := range test.Options {
				// testdata/TestCommonmark/name/input.html
				p := path.Join(t.Name(), "input.html")

				// get the input html from a file
				input, err := ioutil.ReadFile(path.Join("testdata", p))
				if err != nil {
					t.Error(err)
				}

				conv := md.NewConverter("", true, options)
				for _, plugin := range test.Plugins {
					conv.Use(plugin)
				}
				markdown, err := conv.ConvertBytes(input)
				if err != nil {
					t.Error(err)
				}

				// testdata/TestCommonmark/name/output.default.golden
				p = path.Join(t.Name(), "output."+key)
				g.Assert(t, p, markdown)

				gold := goldmark.New(goldmark.WithExtensions(extension.GFM))
				var buf bytes.Buffer
				if err := gold.Convert(markdown, &buf); err != nil {
					t.Error(err)
				}

				// testdata/TestCommonmark/name/goldmark.golden
				p = path.Join(t.Name(), "goldmark")
				g.Assert(t, p, buf.Bytes())
			}
		})
	}
}

func TestCommonmark(t *testing.T) {
	var tests = []GoldenTest{
		{
			Name: "h1",
			Options: map[string]*md.Options{
				"setext": {HeadingStyle: "setext"},
				"atx":    {HeadingStyle: "atx"},
			},
		},
		{
			Name: "h2",
			Options: map[string]*md.Options{
				"setext": {HeadingStyle: "setext"},
				"atx":    {HeadingStyle: "atx"},
			},
		},
		{
			Name: "h3",
			Options: map[string]*md.Options{
				"setext": {HeadingStyle: "setext"},
				"atx":    {HeadingStyle: "atx"},
			},
		},
		{
			Name: "p with content",
		},
		{
			Name: "p inside div",
		},
		{
			Name: "p with span",
		},
		{
			Name: "p with strong",
			Options: map[string]*md.Options{
				"default":    {StrongDelimiter: ""},
				"underscore": {StrongDelimiter: "__"},
			},
		},
		{
			Name: "p with b",
		},
	}

	RunGoldenTest(t, tests)

	/*

		TestCommonmark
		TestPlugins
		TestRules/Keep/Remove

		---- always start with the main tag


		p with content

		strong nested

		h3

		escape


		---- files
		h1.input

		h1.setext.golden
		h1.atx.golden


	*/
}
