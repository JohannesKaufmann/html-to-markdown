package md_test

import (
	"testing"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

func TestPlugins(t *testing.T) {
	var tests = []GoldenTest{
		{
			Name: "strikethrough",
			Variations: map[string]Variation{
				"default": {
					Plugins: []md.Plugin{
						plugin.Strikethrough(""),
					},
				},
			},
		},
		{
			Name: "checkbox",
			Variations: map[string]Variation{
				"default": {
					Plugins: []md.Plugin{
						plugin.TaskListItems(),
					},
				},
			},
		},
		{
			Name:            "table",
			DisableGoldmark: true,
			Variations: map[string]Variation{
				"default": {},
				"table": {
					Plugins: []md.Plugin{
						plugin.Table(),
					},
				},
				"tablecompat": {
					Plugins: []md.Plugin{
						plugin.TableCompat(),
					},
				},
			},
		},
		{
			Name: "movefrontmatter/simple",
			Variations: map[string]Variation{
				"default": {
					Plugins: []md.Plugin{
						plugin.EXPERIMENTALMoveFrontMatter(),
					},
				},
			},
		},
		{
			Name: "movefrontmatter/not",
			Variations: map[string]Variation{
				"default": {
					Plugins: []md.Plugin{
						plugin.EXPERIMENTALMoveFrontMatter(),
					},
				},
			},
		},
		{
			Name: "movefrontmatter/jekyll",
			Variations: map[string]Variation{
				"default": {
					Plugins: []md.Plugin{
						plugin.EXPERIMENTALMoveFrontMatter('-', '+'),
					},
				},
			},
		},
		{
			Name: "movefrontmatter/blog",
			Variations: map[string]Variation{
				"default": {
					Plugins: []md.Plugin{
						plugin.EXPERIMENTALMoveFrontMatter(),
					},
				},
			},
		},
	}

	RunGoldenTest(t, tests)
}
