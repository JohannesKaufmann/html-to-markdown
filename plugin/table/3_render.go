package table

import (
	"strings"
	"unicode/utf8"

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
	counts := calculateMaxCounts(table.Rows)

	// Sometimes a row contains less cells that another row.
	// We then fill it up with empty cells (e.g. "| text |     |").
	table.Rows = fillUpRows(table.Rows, len(counts))

	// - - - - - - - - - - - - - - - - - - - - - - - - - - //

	w.WriteString("\n\n")
	// - - - Header - - - //
	p.writeRow(w, counts, table.Rows[0])
	w.WriteString("\n")
	p.writeHeaderUnderline(w, table.Alignments, counts)
	w.WriteString("\n")

	// - - - Body - - - //
	for _, cells := range table.Rows[1:] {
		p.writeRow(w, counts, cells)
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
func (s *tablePlugin) writeHeaderUnderline(w converter.Writer, alignments []string, counts []int) {
	for i, maxLength := range counts {
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

		if s.padColumns == PadColumnsBehaviorOn {
			w.WriteString(strings.Repeat("-", maxLength))
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

func (s *tablePlugin) writeRow(w converter.Writer, counts []int, cells [][]byte) {
	for i, cell := range cells {
		isFirstCell := i == 0
		if isFirstCell {
			w.WriteString("|")
		}

		currentCount := utf8.RuneCount(cell)
		filler := counts[i] - currentCount

		if s.padColumns == PadColumnsBehaviorOn || s.padColumns == PadColumnsBehaviorSome {
			w.WriteString(" ")
		}

		w.Write(cell)

		if s.padColumns == PadColumnsBehaviorOn && filler > 0 {
			w.WriteString(strings.Repeat(" ", filler))
		}

		if s.padColumns == PadColumnsBehaviorOn || s.padColumns == PadColumnsBehaviorSome {
			w.WriteString(" ")
		}

		w.WriteString("|")
	}
}
