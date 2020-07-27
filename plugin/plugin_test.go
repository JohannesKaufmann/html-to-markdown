package plugin

import (
	"testing"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

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

func TestConfluenceAttachments(t *testing.T) {
	conv := md.NewConverter("", true, nil)
	conv.Use(ConfluenceAttachments())

	input := `<p>Here&rsquo;s an image:</p><p /><ac:image ac:align="center" ac:layout="center" ac:original-height="290" ac:original-width="290"><ri:attachment ri:filename="image.png" ri:version-at-save="1" /></ac:image><p /><p>Another one</p><ac:image ac:align="center" ac:layout="center" ac:original-height="457" ac:original-width="728"><ri:attachment ri:filename="image.jpg" ri:version-at-save="1" /></ac:image><p />`
	expected := `Here’s an image:

![][image.png]

Another one

![][image.jpg]`
	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}

	if markdown != expected {
		t.Errorf("got '%s' but wanted '%s'", markdown, expected)
	}
}
