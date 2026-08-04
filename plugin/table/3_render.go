package table

import (
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

func (p *tablePlugin) renderTable(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	table := p.collectTableContent(ctx, n)
	if table == nil {
		// Sometime we just cannot render the table.
		// Either because it is an empty table OR
		// because there are newlines inside the content (which would break the table).
		return converter.RenderTryNext
	}

	// Sometimes we pad the cells with extra spaces (e.g. "| text    |").
	// For that we first need to know the maximum width of every column.
	widths := calculateMaxWidths(table.Rows)

	// Sometimes a row contains less cells that another row.
	// We then fill it up with empty cells (e.g. "| text |     |").
	table.Rows = fillUpRows(table.Rows, len(widths))

	// - - - - - - - - - - - - - - - - - - - - - - - - - - //

	w.WriteString("\n\n")
	// - - - Header - - - //
	p.writeRow(w, widths, table.Rows[0])
	w.WriteString("\n")
	p.writeHeaderUnderline(w, table.Alignments, widths)
	w.WriteString("\n")

	// - - - Body - - - //
	for _, cells := range table.Rows[1:] {
		p.writeRow(w, widths, cells)
		w.WriteString("\n")
	}

	// - - - Caption - - - //
	if table.Caption != nil {
		w.WriteString("\n\n")
		w.Write(table.Caption)

	}
	// - - - - - - //
	w.WriteString("\n\n")

	return converter.RenderSuccess
}

func getAlignmentFor(alignments []string, index int) string {
	if index > len(alignments)-1 {
		return ""
	}

	return alignments[index]
}
func (s *tablePlugin) writeHeaderUnderline(w converter.Writer, alignments []string, widths []int) {
	for i, maxWidth := range widths {
		align := getAlignmentFor(alignments, i)

		isFirstCell := i == 0
		if isFirstCell {
			w.WriteString("|")
		}
		if align == "left" || align == "center" {
			w.WriteString(":")
		} else {
			w.WriteString("-")
		}

		if s.cellPaddingBehavior == CellPaddingBehaviorAligned {
			w.WriteString(strings.Repeat("-", maxWidth))
		} else {
			w.WriteString("-")
		}

		if align == "right" || align == "center" {
			w.WriteString(":")
		} else {
			w.WriteString("-")
		}
		w.WriteString("|")
	}
}

func (s *tablePlugin) writeRow(w converter.Writer, widths []int, cells [][]byte) {
	for i, cell := range cells {
		isFirstCell := i == 0
		if isFirstCell {
			w.WriteString("|")
		}

		currentWidth := cellDisplayWidth(cell)
		filler := widths[i] - currentWidth

		if s.cellPaddingBehavior == CellPaddingBehaviorAligned || s.cellPaddingBehavior == CellPaddingBehaviorMinimal {
			w.WriteString(" ")
		}

		w.Write(cell)

		if s.cellPaddingBehavior == CellPaddingBehaviorAligned && filler > 0 {
			w.WriteString(strings.Repeat(" ", filler))
		}

		if s.cellPaddingBehavior == CellPaddingBehaviorAligned || s.cellPaddingBehavior == CellPaddingBehaviorMinimal {
			w.WriteString(" ")
		}

		w.WriteString("|")
	}
}
