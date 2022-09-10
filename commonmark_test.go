package md_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Variation struct {
	Options *md.Options
	Plugins []md.Plugin
}
type GoldenTest struct {
	Name   string
	Domain string

	DisableGoldmark bool
	Variations      map[string]Variation
}

func runGoldenTest(t *testing.T, test GoldenTest, variationKey string) {
	variation := test.Variations[variationKey]

	g := goldie.New(t)

	// testdata/TestCommonmark/name/input.html
	p := path.Join(t.Name(), "input.html")

	// get the input html from a file
	input, err := ioutil.ReadFile(path.Join("testdata", p))
	if err != nil {
		t.Error(err)
		return
	}

	if test.Domain == "" {
		test.Domain = "example.com"
	}

	conv := md.NewConverter(test.Domain, true, variation.Options)
	conv.Keep("keep-tag").Remove("remove-tag")
	for _, plugin := range variation.Plugins {
		conv.Use(plugin)
	}
	markdown, err := conv.ConvertBytes(input)
	if err != nil {
		t.Error(err)
	}

	// testdata/TestCommonmark/name/output.default.golden
	p = path.Join(t.Name(), "output."+variationKey)
	g.Assert(t, p, markdown)

	gold := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var buf bytes.Buffer
	if err := gold.Convert(markdown, &buf); err != nil {
		t.Error(err)
	}

	if !test.DisableGoldmark {
		// testdata/TestCommonmark/name/goldmark.golden
		p = path.Join(t.Name(), "goldmark")
		g.Assert(t, p, buf.Bytes())
	}
}

func RunGoldenTest(t *testing.T, tests []GoldenTest) {
	// loop through all test cases that were added manually
	dirs := make(map[string]struct{})
	for _, test := range tests {
		name := test.Name
		name = strings.Replace(name, " ", "_", -1)
		dirs[name] = struct{}{}
	}

	// now add all tests that were found on disk to the tests slice
	err := filepath.Walk(path.Join("testdata", t.Name()),
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}

			// skip folders that don't contain an input.html file
			if _, err := os.Stat(path.Join(p, "input.html")); os.IsNotExist(err) {
				return nil
			}

			parts := strings.SplitN(p, string(os.PathSeparator), 3)
			p = parts[2] // remove "testdata/TestCommonmark/" from "testdata/TestCommonmark/..."

			_, ok := dirs[p]
			if ok {
				return nil
			}

			// add the folder from disk to the tests slice, since its not it there yet
			tests = append(tests, GoldenTest{
				Name: p,
			})
			return nil
		})
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		if len(test.Variations) == 0 {
			test.Variations = map[string]Variation{
				"default": {},
			}
		}

		t.Run(test.Name, func(t *testing.T) {
			if strings.Contains(t.Name(), "#") {
				fmt.Println("the name", test.Name, t.Name(), "seems too be used for multiple tests")
				return
			}

			for variationKey := range test.Variations {
				runGoldenTest(t, test, variationKey)
			}
		})
	}
}

func TestCommonmark(t *testing.T) {
	var tests = []GoldenTest{
		{
			Name:            "link",
			DisableGoldmark: true,
			Variations: map[string]Variation{
				"relative": {
					Options: &md.Options{
						GetAbsoluteURL: func(selec *goquery.Selection, rawURL string, domain string) string {
							return rawURL
						},
					},
				},

				"inlined": {
					Options: &md.Options{LinkStyle: "inlined"},
				},
				"referenced_full": {
					Options: &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "full"},
				},
				"referenced_collapsed": {
					Options: &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "collapsed"},
				},
				"referenced_shortcut": {
					Options: &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "shortcut"},
				},
			},
		},
		{
			Name: "heading",
			Variations: map[string]Variation{
				"atx": {
					Options: &md.Options{HeadingStyle: "atx"},
				},
				"setext": {
					Options: &md.Options{HeadingStyle: "setext"},
				},
			},
		},
		{
			Name: "italic",
			Variations: map[string]Variation{
				"asterisks": {
					Options: &md.Options{EmDelimiter: "*"},
				},
				"underscores": {
					Options: &md.Options{EmDelimiter: "_"},
				},
			},
		},
		{
			Name: "bold",
			Variations: map[string]Variation{
				"asterisks": {
					Options: &md.Options{StrongDelimiter: "**"},
				},
				"underscores": {
					Options: &md.Options{StrongDelimiter: "__"},
				},
			},
		},
		{
			Name: "pre_code",
			Variations: map[string]Variation{
				"indented": {
					Options: &md.Options{CodeBlockStyle: "indented"},
				},
				"fenced_backtick": {
					Options: &md.Options{CodeBlockStyle: "fenced", Fence: "```"},
				},
				"fenced_tilde": {
					Options: &md.Options{CodeBlockStyle: "fenced", Fence: "~~~"},
				},
			},
		},
		{
			Name: "list",
			Variations: map[string]Variation{
				"asterisks": {
					Options: &md.Options{BulletListMarker: "*"},
				},
				"dash": {
					Options: &md.Options{BulletListMarker: "-"},
				},
				"plus": {
					Options: &md.Options{BulletListMarker: "+"},
				},
			},
		},
		{
			Name:            "list_nested",
			DisableGoldmark: true,
			Variations: map[string]Variation{
				"asterisks": {
					Options: &md.Options{BulletListMarker: "*"},
				},
				"dash": {
					Options: &md.Options{BulletListMarker: "-"},
				},
				"plus": {
					Options: &md.Options{BulletListMarker: "+"},
				},
			},
		},
		// + all the test on disk that are added automatically
	}

	RunGoldenTest(t, tests)
}

func TestRealWorld(t *testing.T) {
	var tests = []GoldenTest{
		{
			Name:   "blog.golang.org",
			Domain: "blog.golang.org",
			Variations: map[string]Variation{
				"inlined": {
					Options: &md.Options{LinkStyle: "inlined"},
				},
				"referenced_full": {
					Options: &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "full"},
				},
				"referenced_collapsed": {
					Options: &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "collapsed"},
				},
				"referenced_shortcut": {
					Options: &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "shortcut"},
				},

				"emphasis_asterisks": {
					Options: &md.Options{EmDelimiter: "*", StrongDelimiter: "**"},
				},
				"emphasis_underscores": {
					Options: &md.Options{EmDelimiter: "_", StrongDelimiter: "__"},
				},
			},
		},
		{
			Name:   "golang.org",
			Domain: "golang.org",
		},
		// + all the test on disk that are added automatically
	}
	RunGoldenTest(t, tests)
}
