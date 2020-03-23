package plugin

import (
	"testing"

	md "github.com/JohannesKaufmann/html-to-markdown"
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

func TestTable_simple(t *testing.T) {
	conv := md.NewConverter("", true, nil)
	conv.AddRules(EXPERIMENTAL_Table...)

	input := `
<table style="width:100%">
	<tr>
		<th>Firstname</th>
		<th>Lastname</th>
		<th>Age</th>
	</tr>
	<tr>
		<td>Jill</td>
		<td>Smith</td>
		<td>50</td>
	</tr>
	<tr>
		<td>Eve</td>
		<td>Jackson</td>
		<td>94</td>
	</tr>
</table> 
	`
	expected := `| Firstname | Lastname | Age |
| --- | --- | --- |
| Jill | Smith | 50 |
| Eve | Jackson | 94 |`
	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if markdown != expected {
		t.Errorf("got '%s' but wanted '%s'", markdown, expected)
	}
}

func TestTable_escape_pipe(t *testing.T) {
	conv := md.NewConverter("", true, nil)
	conv.AddRules(EXPERIMENTAL_Table...)

	input := `
<table style="width:100%">
	<tr>
		<th>Firstname</th>
		<th>With | Character</th>
		<th>Age</th>
	</tr>
	<tr>
		<td>Jill</td>
		<td>Smith</td>
		<td>50</td>
	</tr>
	<tr>
		<td>Eve</td>
		<td>Jackson</td>
		<td>94</td>
	</tr>
</table> 
	`
	expected := `| Firstname | With \| Character | Age |
| --- | --- | --- |
| Jill | Smith | 50 |
| Eve | Jackson | 94 |`
	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if markdown != expected {
		t.Errorf("got '%s' but wanted '%s'", markdown, expected)
	}
}
