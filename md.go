package md

import (
	"fmt"
	"strings"
)

func isTextElement(e *Element) bool {
	if e == nil {
		return false
	}
	return e.Tag == TextNode || e.Tag == Bold || e.Tag == Italics || e.Tag == Strikethrough || e.Tag == Link
}

var LineWidth = 80

func wordWrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped
}

func toMD(e *Element, before *Element, after *Element, parent Tag) string {
	switch e.Tag {
	case Text:
		text := ElemToMD(e.Tag, e.ChildNodes)
		text = wordWrap(text, LineWidth)

		if after != nil && (after.Tag == Text || after.Tag == ListItem) {
			text += "\n"
		}

		return text + "\n"
	case TextNode:
		// if after != nil && after.Tag == TextNode {
		// 	e.Text += " "
		// }
		// if isTextElement(after) {
		// 	e.Text += " "
		// }
		return e.Text

	case Header:
		return fmt.Sprintf("\n%s %s\n\n", strings.Repeat("#", e.Level), ElemToMD(e.Tag, e.ChildNodes))
	case Italics:
		text := ElemToMD(e.Tag, e.ChildNodes)
		// text = cleanString(text)
		// if len(text) == 0 {
		// 	return ""
		// }
		text = fmt.Sprintf("_%s_", text)
		if isTextElement(after) {
			text += " "
		}
		return text
	case Bold:
		text := ElemToMD(e.Tag, e.ChildNodes)
		// text = cleanString(text)
		// if len(text) == 0 {
		// 	return ""
		// }
		text = fmt.Sprintf("**%s**", text)
		// if isTextElement(after) {
		// 	text += " "
		// }
		return text
	case Strikethrough:
		text := ElemToMD(e.Tag, e.ChildNodes)
		// text = cleanString(text)
		// if len(text) == 0 {
		// 	return ""
		// }
		text = fmt.Sprintf("~~%s~~", text)
		// if isTextElement(after) {
		// 	text += " "
		// }
		return text

	case ListItem:
		content := fmt.Sprintf("- %s\n", ElemToMD(e.Tag, e.ChildNodes))
		if after == nil || after.Tag != ListItem {
			content += "\n"
		}
		return content
	case Link:
		text := ElemToMD(e.Tag, e.ChildNodes)
		text = strings.TrimSpace(text)
		return fmt.Sprintf("[%s](%s)", text, e.Href)
	case Image:
		return fmt.Sprintf("![%s](%s)", e.Alt, e.Src)
	case Blockquote:
		var texts []string
		for _, text := range strings.Split(ElemToMD(e.Tag, e.ChildNodes), "\n") {
			t := fmt.Sprintf("> %s\n", text)
			texts = append(texts, t)
		}
		// text := fmt.Sprintf("> %s\n", ElemToMD(e.Tag, e.ChildNodes))
		// if after != nil && after.Tag != Blockquote {
		// 	text += "\n"
		// }

		return strings.Join(texts, "") + "\n"
	case Divider:
		return "\n---\n\n"
	case InlineCode:
		return fmt.Sprintf("`%s`", e.Text)
	case CodeBlock:
		return fmt.Sprintf("```\n%s\n```", e.Text)

	case "BREAK":
		if parent != Italics && parent != Bold && parent != Strikethrough {
			return "\n"
		}

	default:
		fmt.Println("cant convert element to markdown", e.Tag)

	}
	return ""
}

func ElemToMD(parent Tag, elements []*Element) string {
	// type Item struct {
	// 	Before MDTag
	// 	Tag    Tag
	// 	Value  string
	// 	After  X
	// }
	// for loop
	// 		which one has precedence?
	// wbr: "None" has precedence if surrounding is text_nodes
	//				if next is other p tag -> "NewLine"
	var builder strings.Builder

	for i, element := range elements {
		var before *Element
		var after *Element
		if i != 0 {
			before = elements[i-1]
		}
		if i < len(elements)-1 {
			after = elements[i+1]
		}
		builder.WriteString(
			toMD(element, before, after, parent),
		)
	}
	res := builder.String()

	res = strings.TrimSpace(res)
	// remove all unnecessary new line characters
	res = newLineRegex.ReplaceAllString(res, "\n\n")

	return res
}
