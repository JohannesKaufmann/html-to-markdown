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
			Name: "table/simple",
			Plugins: []md.Plugin{
				func(conv *md.Converter) []md.Rule {
					return plugin.EXPERIMENTAL_Table
				},
			},
		},
		{
			Name: "table/escape pipe",
			Plugins: []md.Plugin{
				func(conv *md.Converter) []md.Rule {
					return plugin.EXPERIMENTAL_Table
				},
			},
		},
	}

	RunGoldenTest(t, tests)
}
