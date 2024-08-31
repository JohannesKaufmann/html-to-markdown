package tester

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JohannesKaufmann/dom"
	"golang.org/x/net/html"
)

func Parse(t *testing.T, rawHTML string, startFrom string) *html.Node {
	if startFrom == "" {
		startFrom = "body"
	}

	rawHTML = strings.TrimSpace(rawHTML)

	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	err = html.Render(&b, doc)
	if err != nil {
		t.Error(err)
	}

	var foundNode *html.Node
	var finder func(node *html.Node)
	finder = func(node *html.Node) {
		if foundNode != nil {
			return
		}
		if dom.NodeName(node) == startFrom {
			foundNode = node
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			finder(child)
		}
	}
	finder(doc)

	if foundNode == nil {
		t.Error("could not find node for 'startFrom'")
	}
	return foundNode
}
