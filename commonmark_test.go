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
	"github.com/sebdah/goldie/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type GoldenTest struct {
	Name   string
	Domain string

	Options map[string]*md.Options
	Plugins []md.Plugin
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
		if len(test.Options) == 0 {
			test.Options = map[string]*md.Options{
				"default": nil,
			}
		}

		t.Run(test.Name, func(t *testing.T) {
			if strings.Contains(t.Name(), "#") {
				fmt.Println("the name", test.Name, t.Name(), "seems too be used for multiple tests")
				return
			}

			g := goldie.New(t)

			for key, options := range test.Options {
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

				conv := md.NewConverter(test.Domain, true, options)
				conv.Keep("keep-tag").Remove("remove-tag")
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
			Name: "link",
			Options: map[string]*md.Options{
				"inlined":              &md.Options{LinkStyle: "inlined"},
				"referenced_full":      &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "full"},
				"referenced_collapsed": &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "collapsed"},
				"referenced_shortcut":  &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "shortcut"},
			},
		},
		{
			Name: "heading",
			Options: map[string]*md.Options{
				"atx":    &md.Options{HeadingStyle: "atx"},
				"setext": &md.Options{HeadingStyle: "setext"},
			},
		},
		{
			Name: "italic",
			Options: map[string]*md.Options{
				"asterisks":   &md.Options{EmDelimiter: "*"},
				"underscores": &md.Options{EmDelimiter: "_"},
			},
		},
		{
			Name: "bold",
			Options: map[string]*md.Options{
				"asterisks":   &md.Options{StrongDelimiter: "**"},
				"underscores": &md.Options{StrongDelimiter: "__"},
			},
		},
		{
			Name: "pre_code",
			Options: map[string]*md.Options{
				"indented":        &md.Options{CodeBlockStyle: "indented"},
				"fenced_backtick": &md.Options{CodeBlockStyle: "fenced", Fence: "```"},
				"fenced_tilde":    &md.Options{CodeBlockStyle: "fenced", Fence: "~~~"},
			},
		},
		{
			Name: "list",
			Options: map[string]*md.Options{
				"asterisks": &md.Options{BulletListMarker: "*"},
				"dash":      &md.Options{BulletListMarker: "-"},
				"plus":      &md.Options{BulletListMarker: "+"},
			},
		},
		{
			Name: "list_nested",
			Options: map[string]*md.Options{
				"asterisks": &md.Options{BulletListMarker: "*"},
				"dash":      &md.Options{BulletListMarker: "-"},
				"plus":      &md.Options{BulletListMarker: "+"},
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
			Options: map[string]*md.Options{
				"inlined":              &md.Options{LinkStyle: "inlined"},
				"referenced_full":      &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "full"},
				"referenced_collapsed": &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "collapsed"},
				"referenced_shortcut":  &md.Options{LinkStyle: "referenced", LinkReferenceStyle: "shortcut"},
				"emphasis_asterisks":   &md.Options{EmDelimiter: "*", StrongDelimiter: "**"},
				"emphasis_underscores": &md.Options{EmDelimiter: "_", StrongDelimiter: "__"},
			},
		},
		{
			Name:   "golang.org",
			Domain: "golang.org",
		},
		{
			Name:   "bonnerruderverein.de",
			Domain: "bonnerruderverein.de",
		},
		// + all the test on disk that are added automatically
	}
	RunGoldenTest(t, tests)
}
