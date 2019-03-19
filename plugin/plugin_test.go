package plugin

import (
	"testing"

	"github.com/mmelvin0/html-to-markdown"
)

func TestStrikethroughDefault(t *testing.T) {
	conv := md.NewConverter("", true, nil)
	conv.Use(Strikethrough(""))

	input := `<p>Strikethrough uses two tildes. <del>Scratch this.</del></p>`
	expected := `Strikethrough uses two tildes. ~Scratch this.~`
	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if markdown != expected {
		t.Errorf("got '%s' but wanted '%s'", markdown, expected)
	}
}
func TestStrikethrough(t *testing.T) {
	conv := md.NewConverter("", true, nil)
	conv.Use(Strikethrough("~~"))

	input := `<p>Strikethrough uses two tildes. <del>Scratch this.</del></p>`
	expected := `Strikethrough uses two tildes. ~~Scratch this.~~`
	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if markdown != expected {
		t.Errorf("got '%s' but wanted '%s'", markdown, expected)
	}
}

func TestTaskListItems(t *testing.T) {
	conv := md.NewConverter("", true, nil)
	conv.Use(TaskListItems())

	input := `
	<ul>
		<li><input type=checkbox checked>Checked!</li>
		<li><input type=checkbox>Check Me!</li>
	</ul>
	`
	expected := `- [x] Checked!
- [ ] Check Me!`
	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if markdown != expected {
		t.Errorf("got '%s' but wanted '%s'", markdown, expected)
	}
}
