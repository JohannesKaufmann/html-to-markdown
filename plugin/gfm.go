package plugin

import (
	"github.com/JohannesKaufmann/html-to-markdown"
)

// GFM is GitHub Flavored Markdown
var GFM []md.Rule

func init() {
	GFM = append(GFM, Strikethrough...)
	GFM = append(GFM, Table...)
	GFM = append(GFM, TaskListItems...)
}
