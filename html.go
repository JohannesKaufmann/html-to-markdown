package md

import (
	"regexp"
	"strings"
)

type Tag string

var (
	Root Tag = "ROOT"

	Text     Tag = "TEXT"
	TextNode Tag = "TEXT_NODE"

	Header        Tag = "HEADER"          // # text
	Italics       Tag = "ITALICS"         // *text*
	Bold          Tag = "BOLD"            // __text__
	Strikethrough Tag = "STRIKE_THROUGH"  // ~~text~~
	ListItem      Tag = "LIST_ITEM"       // - text
	Link          Tag = "LINK"            // [text](href)
	Image         Tag = "IMAGE"           // ![alt](src)
	Blockquote    Tag = "BLOCK_QUOTE"     // > text
	Divider       Tag = "HORIZONTAL_RULE" // ---

	CodeBlock  Tag = "CODE_BLOCK"  // starts with ```js
	InlineCode Tag = "INLINE_CODE" // `variable`

	// Tables
	// Br
)

type Element struct {
	Tag   Tag
	Level int `json:",omitempty"`

	// Text should ONLY be used by a text node
	Text string `json:",omitempty"`

	// a tag
	Href string `json:",omitempty"`

	// img tag
	Src string `json:",omitempty"`
	Alt string `json:",omitempty"`

	ChildNodes []*Element

	// TODO: only for rare cases
	Data interface{}
}

var spaceRegex = regexp.MustCompile(`\s+`)
var tabRegex = regexp.MustCompile(`\\t+`)
var newLineRegex = regexp.MustCompile(`[\n]{2,}`)

func removeTabs(text string) string {
	return strings.Map(func(r rune) rune {
		if r == '\t' { // || r == '\n'
			return -1
		}
		return r
	}, text)
}

// func cleanString(text string) string {
// 	text = strings.TrimSpace(text)
// 	text = spaceRegex.ReplaceAllString(text, " ")
// 	return text
// }
func isInlineNode(t Tag) bool {
	switch t {
	case TextNode, Italics, Bold, Strikethrough, Link, Image:
		return true
	}
	return false
}

// func clean(text string) string {
// 	return strings.Map(func(r rune) rune {
// 		if unicode.IsSpace(r) {
// 			return -1
// 		}
// 		return r
// 	}, text)
// }

/*
func _SelecToElem(domain string, isChildren bool, selec *goquery.Selection) []*Element {
	var elements []*Element

	selec.Contents().Each(func(i int, s *goquery.Selection) {
		node := goquery.NodeName(s)
		e := &Element{}

		switch node {
		// - - special cases - - //
		case "figure":
			e.Tag = Image
			e.Src = s.Find("img").AttrOr("src", "")
			e.Alt = s.Find("figcaption").Text()

			u, err := url.Parse(e.Src)
			if err != nil {
				log.Fatal(err)
			}
			if !u.IsAbs() {
				e.Src = domain + e.Src
			}

			elements = append(elements, e)
			// TODO: look for <figcaption> and add it to the image if found
			// elements = append(elements, getElements(domain, isChildren, s)...)

		case "iframe":
			e.Tag = Link
			e.Href = s.AttrOr("src", "")

			if strings.Contains(e.Href, "facebook.com") {
				fmt.Println("facebook iframe")
			} else if strings.Contains(e.Href, "twitter.com") {
				fmt.Println("twitter iframe", s.AttrOr("title", "NO_TITLE"))
			} else if strings.Contains(e.Href, "youtube.com") {
				fmt.Println("youtube iframe")
				// https://www.youtube.com/embed/65gN8xbM1BY?showinfo=0&rel=0&iv_load_policy=3
				// [![IMAGE ALT TEXT HERE](http://img.youtube.com/vi/YOUTUBE_VIDEO_ID_HERE/0.jpg)](http://www.youtube.com/watch?v=YOUTUBE_VIDEO_ID_HERE)

				// a: http://www.youtube.com/watch?feature=player_embedded&v=YOUTUBE_VIDEO_ID_HERE
				// img: http://img.youtube.com/vi/YOUTUBE_VIDEO_ID_HERE/0.jpg

			} else {
				fmt.Println("other iframe:", e.Href)
			}
			elements = append(elements, e)

		case "audio":
			e.Tag = Link
			e.Href = s.AttrOr("src", "")
			elements = append(elements, e)
		case "time":
			elements = append(elements, SelecToElem(domain, isChildren, s)...)
			// TODO: parse <time class="entry-date" datetime="2018-04-26T18:17:14+00:00">26. April 2018</time>

		case "wbr": // A text with word break opportunities
			// START<wbr>text</wbr>END -> "STARTtextEND" without spaces
			elements = append(elements, SelecToElem(domain, isChildren, s)...)
			// TODO: without spaces

		// - - normal cases - - //
		case "span":
			// TODO: what about span?
			elements = append(elements, SelecToElem(domain, isChildren, s)...)

		case "p":
			e.Tag = Text
			e.ChildNodes = SelecToElem(domain, true, s)
			if len(e.ChildNodes) != 0 {
				elements = append(elements, e)
			}

		case "#text":
			e.Tag = TextNode
			e.Text = s.Text()
			e.Text = removeTabs(e.Text)

			trimed := strings.TrimSpace(e.Text)
			// e.Text = cleanString(s.Text())

			if trimed != "" {
				elements = append(elements, e)
			}

		case "h1", "h2", "h3", "h4", "h5", "h6":
			e.Tag = Header
			l, err := strconv.Atoi(node[1:])
			if err != nil {
				log.Fatal(err)
			}
			e.Level = l

			e.ChildNodes = SelecToElem(domain, true, s)
			elements = append(elements, e)

		case "cite", "em", "i":
			e.Tag = Italics
			e.ChildNodes = SelecToElem(domain, true, s)
			if len(e.ChildNodes) != 0 {
				elements = append(elements, e)
			}

		case "strong", "b":
			e.Tag = Bold
			e.ChildNodes = SelecToElem(domain, true, s)
			if len(e.ChildNodes) != 0 {
				elements = append(elements, e)
			}
		case "del":
			e.Tag = Strikethrough
			e.ChildNodes = SelecToElem(domain, true, s)
			if len(e.ChildNodes) != 0 {
				elements = append(elements, e)
			}
		case "li":
			e.Tag = ListItem
			e.ChildNodes = SelecToElem(domain, true, s)
			if len(e.ChildNodes) != 0 {
				elements = append(elements, e)
			}

		case "a":
			e.Tag = Link
			e.Href = s.AttrOr("href", "")
			if strings.HasPrefix(e.Href, "javascript:") {
				e.Href = ""
			}
			e.ChildNodes = SelecToElem(domain, true, s)

			if len(e.ChildNodes) != 0 {
				elements = append(elements, e)
			}

		case "picture":
			s = s.Find("img")
			fallthrough
		case "img":
			e.Tag = Image
			e.Src = s.AttrOr("src", "")
			u, err := url.Parse(e.Src)
			if err != nil {
				log.Fatal(err)
			}
			if !u.IsAbs() {
				e.Src = domain + e.Src
			}

			e.Alt = s.AttrOr("alt", "")
			elements = append(elements, e)

		case "blockquote":
			e.Tag = Blockquote
			e.ChildNodes = SelecToElem(domain, true, s)
			elements = append(elements, e)

		case "hr":
			e.Tag = Divider
			elements = append(elements, e)

		case "code":
			e.Tag = InlineCode
			e.Text = s.Text()
			elements = append(elements, e)

		case "pre":
			e.Tag = CodeBlock
			e.Text = s.Text()
			e.Text = strings.TrimSpace(e.Text)
			elements = append(elements, e)

		case "div",
			"ul",
			"ol",
			"section",
			"article",
			"aside",
			"footer",
			"nav",
			"header",
			"body",
			"html":
			elems := SelecToElem(domain, isChildren, s)
			if len(elems) != 0 {
				var onlyTextNodes = true
				for _, elem := range elems {
					if !isInlineNode(elem.Tag) {
						onlyTextNodes = false
					}
				}
				if onlyTextNodes {
					e.Tag = Text
					e.ChildNodes = elems
					elements = append(elements, e)
				} else {
					elements = append(elements, elems...)
				}
			}

		case "br":
			e.Tag = "BREAK"
			elements = append(elements, e)
			// ignore for now?

		case "script", "style", "#comment", "head", "svg":
			// ignore
		default:
			fmt.Println("[get parent nodes] unknown tag", node, e.Tag)
			elements = append(elements, SelecToElem(domain, isChildren, s)...)
		}
	})
	return elements
}
*/
