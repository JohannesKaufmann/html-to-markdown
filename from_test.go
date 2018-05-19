package md

import (
	"fmt"
	"testing"
)

func TestFromString(t *testing.T) {
	// "<p>Some Text</p>"
	markdown, err := FromString("", `<ul>
		<li>Some Thing</li>
		<li>Another Thing</li>
	</ul>

	<p>Some Text</p>
	
	<ol>
		<li>First Thing</li>
		<li>Second Thing</li>
	</ol>
	`)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("- - - result - - -")
	fmt.Println(markdown)

	t.Fail()
}
