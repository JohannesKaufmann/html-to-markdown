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
			Plugins: []md.Plugin{
				plugin.Strikethrough(""),
			},
		},
		{
			Name: "strikethrough with space",
			Plugins: []md.Plugin{
				plugin.Strikethrough(""),
			},
		},
		{ // #23
			Name: "strikethrough next to each other",
			Plugins: []md.Plugin{
				plugin.Strikethrough(""),
			},
		},
		{
			Name: "checkbox",
			Plugins: []md.Plugin{
				plugin.TaskListItems(),
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
	}

	RunGoldenTest(t, tests)
}
