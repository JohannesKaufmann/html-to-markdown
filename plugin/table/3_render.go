package table

import (
	"strings"
	"unicode/utf8"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

func (s *tablePlugin) renderTableBody(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	table := collectTableContent(ctx, n)
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
	s.writeRow(w, counts, table.Rows[0])
	w.WriteString("\n")
	s.writeHeaderUnderline(w, counts)
	w.WriteString("\n")

	// - - - Body - - - //
	for _, cells := range table.Rows[1:] {
		s.writeRow(w, counts, cells)
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

func (s *tablePlugin) writeHeaderUnderline(w converter.Writer, counts []int) {
	for i, maxLength := range counts {
		isFirstCell := i == 0
		if isFirstCell {
			w.WriteString("|")
		}
		w.WriteString(" ")
		w.WriteString(strings.Repeat("-", maxLength))

		// TODO: maybe no spaces? So for example "|----|" instead of "| --- |"
		w.WriteString(" |")
	}
}

func (s *tablePlugin) writeRow(w converter.Writer, counts []int, cells [][]byte) {
	for i, cell := range cells {
		isFirstCell := i == 0
		if isFirstCell {
			w.WriteString("|")
		}
		w.WriteString(" ")
		w.Write(cell)

		currentCount := utf8.RuneCount(cell)
		filler := counts[i] - currentCount

		if filler > 0 {
			w.WriteString(strings.Repeat(" ", filler))
		}

		w.WriteString(" |")
	}
}
