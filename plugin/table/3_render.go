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

	// TODO: find better name / function
	x := make([][][]byte, 0, 1+len(table.BodyRows))
	x = append(x, table.HeaderRow)
	x = append(x, table.BodyRows...)
	counts := calculateMaxCounts(x)

	if len(table.HeaderRow) == 0 {
		// There needs to be *header* row so that the table is recognized.
		// So it is better to have an empty header row...
		var emptyCells [][]byte
		for range counts {
			emptyCells = append(emptyCells, []byte(""))
		}
		table.HeaderRow = emptyCells
	}

	// - - - - - - - - - - - - - - - - - - - - - - - - - - //

	w.WriteString("\n\n")
	// - - - Header - - - //
	s.writeRow(w, counts, table.HeaderRow)
	w.WriteString("\n")
	s.writeHeaderUnderline(w, counts)
	w.WriteString("\n")

	// - - - Body - - - //
	for _, cells := range table.BodyRows {
		s.writeRow(w, counts, cells)
		w.WriteString("\n")
	}
	// - - - - - - //
	w.WriteString("\n\n")

	if table.Caption != nil {
		w.Write(table.Caption)
		w.WriteString("\n\n")
	}

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
