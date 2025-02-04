package table

import (
	"bytes"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type tableContent struct {
	Alignments []string
	Rows       [][][]byte
	Caption    []byte
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
		case "table":
			// This will be caught with the newline check anyway.
			// But we can safe some effort by aborting early...
			return true
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
		Alignments: collectAlignments(headerRowNode, normalRowNodes),
		Rows:       rows,
		Caption:    collectCaption(ctx, node),
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

func getFirstNode(node *html.Node, nodes ...*html.Node) *html.Node {
	if node != nil {
		return node
	}
	if len(nodes) >= 1 {
		return nodes[0]
	}
	return nil
}

func collectAlignments(headerRowNode *html.Node, rowNodes []*html.Node) []string {
	firstRow := getFirstNode(headerRowNode, rowNodes...)
	if firstRow == nil {
		return nil
	}

	cellNodes := dom.FindAllNodes(firstRow, func(node *html.Node) bool {
		name := dom.NodeName(node)
		return name == "th" || name == "td"
	})

	var alignments []string
	for _, cellNode := range cellNodes {
		align := dom.GetAttributeOr(cellNode, "align", "")

		alignments = append(alignments, align)
	}

	return alignments
}
func (p *tablePlugin) collectCellsInRow(ctx converter.Context, rowIndex int, rowNode *html.Node) ([][]byte, []modification) {
	cellNodes := dom.FindAllNodes(rowNode, func(node *html.Node) bool {
		name := dom.NodeName(node)
		return name == "th" || name == "td"
	})

	cellContents := make([][]byte, 0, len(cellNodes))
	modifications := make([]modification, 0)

	for index, cellNode := range cellNodes {
		var buf bytes.Buffer
		ctx.RenderNodes(ctx, &buf, cellNode)

		content := buf.Bytes()
		content = bytes.TrimSpace(content)

		content = ctx.UnEscapeContent(content)

		cellContents = append(cellContents, content)

		// - - col / row span - - //
		rowSpan := getNumberAttributeOr(cellNode, "rowspan", 1)
		colSpan := getNumberAttributeOr(cellNode, "colspan", 1)

		mods := calculateModifications(rowIndex, index, rowSpan, colSpan, p.getContentForMergedCell(content))

		modifications = append(modifications, mods...)
	}

	return cellContents, modifications
}
func (p *tablePlugin) collectRows(ctx converter.Context, headerRowNode *html.Node, rowNodes []*html.Node) [][][]byte {
	rowContents := make([][][]byte, 0, len(rowNodes)+1)
	groupedModifications := make([][]modification, 0)

	// - - 1. the header row - - //
	if headerRowNode != nil {
		cells, mods := p.collectCellsInRow(ctx, 0, headerRowNode)

		rowContents = append(rowContents, cells)
		groupedModifications = append(groupedModifications, mods)
	} else {
		// There needs to be *header* row so that the table is recognized.
		// So it is better to have an empty header row...
		rowContents = append(rowContents, [][]byte{})
	}

	// - - 2. the normal rows - - //
	for index, rowNode := range rowNodes {
		cells, mods := p.collectCellsInRow(ctx, index+1, rowNode)

		rowContents = append(rowContents, cells)
		groupedModifications = append(groupedModifications, mods)
	}

	// Sometimes a cell wants to *span* over multiple columns or/and rows.
	// We collected these modifications and are now applying it,
	// by shifting the cells around.
	rowContents = applyGroupedModifications(rowContents, groupedModifications)

	if p.skipEmptyRows {
		rowContents = removeEmptyRows(rowContents)
	}
	if p.promoteFirstRowToHeader {
		rowContents = removeFirstRowIfEmpty(rowContents)
	}

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
