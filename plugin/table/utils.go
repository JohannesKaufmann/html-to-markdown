package table

import (
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/JohannesKaufmann/dom"
	"golang.org/x/net/html"
)

func calculateMaxCounts(rows [][][]byte) []int {
	maxCounts := make([]int, 0)

	for _, cells := range rows {
		for index, cell := range cells {
			count := utf8.RuneCount(cell)

			if index >= len(maxCounts) {
				maxCounts = append(maxCounts, 0)
			}
			currentMax := maxCounts[index]
			if count > currentMax {
				maxCounts[index] = count
			}
		}
	}
	return maxCounts
}

func fillUpRows(rows [][][]byte, maxColumnCount int) [][][]byte {

	for i, cells := range rows {
		missingCells := maxColumnCount - len(cells)
		for range missingCells {
			rows[i] = append(rows[i], []byte(""))
		}
	}

	return rows
}

func getNumberAttributeOr(node *html.Node, key string, fallback int) int {
	val, ok := dom.GetAttribute(node, key)
	if !ok {
		return fallback
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	if num < 1 {
		return fallback
	}

	return num
}

type modification struct {
	y    int
	x    int
	data []byte
}

func calculateModifications(currentRowIndex, currentColIndex, rowSpan, colSpan int, data []byte) []modification {

	mods := make([]modification, 0)

	if colSpan <= 1 && rowSpan <= 1 {
		// No modification is needed
		return mods
	}

	// Calculate modifications for colspan
	for dx := 1; dx < colSpan; dx++ {
		// Add modifications for the same row
		mods = append(mods, modification{
			y:    currentRowIndex,
			x:    currentColIndex + dx,
			data: data,
		})
	}

	// Calculate modifications for subsequent rows
	if rowSpan > 1 {
		for dy := 1; dy < rowSpan; dy++ {
			for dx := 0; dx < colSpan; dx++ {
				mods = append(mods, modification{
					y:    currentRowIndex + dy,
					x:    currentColIndex + dx,
					data: data,
				})
			}
		}
	}

	return mods
}

// TODO: better name?
func applyModifications(contents [][][]byte, mods []modification) [][][]byte {
	for _, mod := range mods {
		// Grow on the y axis
		contents = growSlice(contents, mod.y, nil)

		// Grow on the x axis
		// (Note: we only grow x-1 since `Insert` takes care of the rest)
		contents[mod.y] = growSlice(contents[mod.y], mod.x-1, nil)

		// Now we can do our change:
		contents[mod.y] = slices.Insert(contents[mod.y], mod.x, mod.data)
	}

	return contents
}

// growSlice ensures the slice has enough capacity to access the given index.
func growSlice[T any](contents []T, index int, placeholderVal T) []T {
	// Calculate the required growth
	currentLen := len(contents)
	if index < currentLen {
		return contents
	}

	growBy := index - currentLen + 1

	// Grow the slice by appending values
	for range growBy {
		contents = append(contents, placeholderVal)
	}

	return contents
}
