package table

import (
	"bytes"
	"strings"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type tableContent struct {
	Rows    [][][]byte
	Caption []byte
}

func containsNewline(b []byte) bool {
	return bytes.Contains(b, []byte("\n"))
}

func containsProblematicNode(node *html.Node) bool {
	problematicNode := dom.FindFirstNode(node, func(n *html.Node) bool {
		name := dom.NodeName(n)

		if dom.NameIsHeading(name) {
			return true
		}
		switch name {
		case "br", "hr", "ul", "ol", "blockquote":
			return true
		}

		return false
	})

	return problematicNode != nil
}

func (p *tablePlugin) collectTableContent(ctx converter.Context, node *html.Node) *tableContent {
	if containsProblematicNode(node) {
		// There are certain nodes (e.g. <hr />) that cannot be in a table.
		// If we found one, we unfortunately cannot convert the table.
		//
		// Note: It is okay for a block node (e.g. <div>) to be in a table.
		//       However once it causes multiple lines, it does not work anymore.
		//       For that we have the `containsNewline` check below.
		return nil
	}

	headerRowNode := selectHeaderRowNode(node)
	normalRowNodes := selectNormalRowNodes(node, headerRowNode)

	rows := p.collectRows(ctx, headerRowNode, normalRowNodes)
	if len(rows) == 0 {
		return nil
	}

	for _, cells := range rows {
		for _, cell := range cells {
			if containsNewline(cell) {
				// Having newlines inside the content would break the table.
				// So unfortunately we cannot convert the table.
				return nil
			}
		}
	}

	return &tableContent{
		Rows:    rows,
		Caption: collectCaption(ctx, node),
	}
}

// Sometimes a cell wants to *span* over multiple columns or/and rows.
// What should be displayed in those other cells?
// Render exactly the same content OR an empty string?
func (p *tablePlugin) getContentForMergedCell(originalContent []byte) []byte {
	if p.mergeContentReplication {
		return originalContent
	}

	return []byte("")
}

func (p *tablePlugin) collectCellsInRow(ctx converter.Context, rowIndex int, rowNode *html.Node) ([][]byte, []modification) {
	if rowNode == nil {
		return nil, nil
	}

	name := dom.NodeName(rowNode)
	if name != "tr" {
		panic("the table child is not a tr but " + name)
	}

	// TODO: we should not use child nodes but instead get directly the td and th
	cellNodes := dom.AllChildNodes(rowNode)
	cellsContents := make([][]byte, 0, len(cellNodes))
	modifications := make([]modification, 0)

	var index int
	for _, cellNode := range cellNodes {
		name := dom.NodeName(cellNode)
		if name == "#text" && strings.TrimSpace(dom.CollectText(cellNode)) == "" {
			continue
		}
		if name != "td" && name != "th" {
			panic("the table subchild is not a td but " + name)
		}

		var buf bytes.Buffer
		ctx.RenderNodes(ctx, &buf, cellNode)

		content := buf.Bytes()
		content = bytes.TrimSpace(content)

		content = ctx.UnEscapeContent(content)

		cellsContents = append(cellsContents, content)

		// - - col / row span - - //
		rowSpan := getNumberAttributeOr(cellNode, "rowspan", 1)
		colSpan := getNumberAttributeOr(cellNode, "colspan", 1)

		mods := calculateModifications(rowIndex, index, rowSpan, colSpan, p.getContentForMergedCell(content))

		modifications = append(modifications, mods...)

		index++
	}

	return cellsContents, modifications
}
func (p *tablePlugin) collectRows(ctx converter.Context, headerRowNode *html.Node, rowNodes []*html.Node) [][][]byte {
	rowContents := make([][][]byte, 0, len(rowNodes)+1)
	modifications := make([]modification, 0)

	// - - 1. the header row - - //
	if headerRowNode != nil {
		cells, mods := p.collectCellsInRow(ctx, 0, headerRowNode)

		rowContents = append(rowContents, cells)
		modifications = append(modifications, mods...)
	} else {
		// There needs to be *header* row so that the table is recognized.
		// So it is better to have an empty header row...
		rowContents = append(rowContents, [][]byte{})
	}

	// - - 2. the normal rows - - //
	for index, rowNode := range rowNodes {
		cells, mods := p.collectCellsInRow(ctx, index+1, rowNode)

		rowContents = append(rowContents, cells)
		modifications = append(modifications, mods...)
	}

	// Sometimes a cell wants to *span* over multiple columns or/and rows.
	// We collected these modifications and are now applying it,
	// by shifting the cells around.
	applyModifications(rowContents, modifications)

	return rowContents
}

func collectCaption(ctx converter.Context, node *html.Node) []byte {
	captionNode := dom.FindFirstNode(node, func(node *html.Node) bool {
		return node.DataAtom == atom.Caption
	})
	if captionNode == nil {
		return nil
	}

	var buf bytes.Buffer
	ctx.RenderNodes(ctx, &buf, captionNode)

	content := buf.Bytes()
	content = bytes.TrimSpace(content)

	return content
}
