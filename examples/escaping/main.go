package main

import (
	"fmt"
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func main() {
	html := `<p>fake **bold** and real <strong>bold</strong></p>`
	// With "basic" we get:
	// "fake \*\*bold\*\* and real **bold**"
	// which would render as:
	// "<p>fake **bold** and real <strong>bold</strong></p>"

	// With "none" we get:
	// "fake **bold** and real **bold**"
	// which would render as:
	// "<p>fake <strong>bold</strong> and real <strong>bold</strong></p>"

	opt := &md.Options{
		EscapeMode: "basic", // default
	}
	conv := md.NewConverter("", true, opt)

	markdown1, err := conv.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

	// - - - - //

	opt = &md.Options{
		EscapeMode: "disabled",
	}
	conv = md.NewConverter("", true, opt)

	markdown2, err := conv.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("with basic:", markdown1)
	fmt.Println("with disabled:", markdown2)
}
