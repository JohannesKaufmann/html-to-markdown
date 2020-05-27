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

func TestConfluenceCodeBlock(t *testing.T) {
	conv := md.NewConverter("", true, nil)
	conv.Use(ConfluenceCodeBlock())

	input := `<ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="150db472-e155-47c7-a195-c581bf891af5"><ac:plain-text-body><![CDATA[FOR stuff IN imdb_vertices
	FILTER LIKE(stuff.description, "%good%vs%evil%", true)
  RETURN stuff.description]]></ac:plain-text-body></ac:structured-macro>
some other stuff
<ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="150db472-e155-47c7-a195-c581bf891af5"><ac:parameter ac:name="language">sql</ac:parameter><ac:plain-text-body><![CDATA[FOR stuff IN imdb_vertices
	FILTER LIKE(stuff.description, "%good%vs%evil%", true)
  RETURN stuff.description]]></ac:plain-text-body></ac:structured-macro>`
	expected := "```" + `
FOR stuff IN imdb_vertices
	FILTER LIKE(stuff.description, "%good%vs%evil%", true)
  RETURN stuff.description
` + "```" + `
some other stuff
` + "```sql" + `
FOR stuff IN imdb_vertices
	FILTER LIKE(stuff.description, "%good%vs%evil%", true)
  RETURN stuff.description
` + "```"
	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if markdown != expected {
		t.Errorf("got '%s' but wanted '%s'", markdown, expected)
	}
}
