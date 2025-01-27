package table

import (
	"bytes"
	"strings"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

type tableContent struct {
	HeaderRow [][]byte
	BodyRows  [][][]byte
}

func containsNewline(b []byte) bool {
	return bytes.Contains(b, []byte("\n"))
}

func collectTableContent(ctx converter.Context, node *html.Node) *tableContent {
	headerRowNode := selectHeaderRowNode(node)
	normalRowNodes := selectNormalRowNodes(node, headerRowNode)

	headerRow, modifications := collectCellsInRow(ctx, 0, headerRowNode)
	// TODO: the modifications should also affect the body content
	applyModifications([][][]byte{headerRow}, modifications)

	bodyRows := collectRows(ctx, normalRowNodes)

	if len(headerRow) == 0 && len(bodyRows) == 0 {
		return nil
	}
	for _, cell := range headerRow {
		if containsNewline(cell) {
			return nil
		}
	}
	for _, cells := range bodyRows {
		for _, cell := range cells {
			if containsNewline(cell) {
				return nil
			}
		}
	}

	return &tableContent{
		HeaderRow: headerRow,
		BodyRows:  bodyRows,
	}
}

func collectCellsInRow(ctx converter.Context, rowIndex int, rowNode *html.Node) ([][]byte, []modification) {
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

		mods := calculateModifications(rowIndex, index, rowSpan, colSpan)

		modifications = append(modifications, mods...)

		index++
	}
	return cellsContents, modifications
}
func collectRows(ctx converter.Context, rowNodes []*html.Node) [][][]byte {
	rowContents := make([][][]byte, 0, len(rowNodes))
	modifications := make([]modification, 0)

	for index, rowNode := range rowNodes {
		cells, mods := collectCellsInRow(ctx, index, rowNode)
		modifications = append(modifications, mods...)

		rowContents = append(rowContents, cells)
	}

	applyModifications(rowContents, modifications)

	return rowContents
}
