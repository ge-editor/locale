package locale

import (
	"github.com/gdamore/tcell/v3"
)

// Cell represents a single character Cell on screen.
type Cell struct {
	Ch                   rune // int32
	Style                tcell.Style
	Size                 int
	Width                int // screen cell width
	Class                CharClass
	TotalWidthLogicalRow int // screen cell total width of logical row without hanging indent
	LogicalRowIndex      int
}

/*
	type Cell struct {
	    Ch rune                  // 4
	    Style tcell.Style        // ?
	    Size uint8               // 1
	    Width uint8              // 1
	    Class locale.CharClass   // 1 or 2
	    LogicalRowIndex uint16   // 2
	    TotalWidthLogicalRow uint16
	}
*/

func (m *Cell) IsEmpty() bool {
	return m.Class == 0
}

func (c *Cell) Clear() {
	*c = Cell{}
}
