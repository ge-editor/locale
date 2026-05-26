package locale

import (
	"github.com/gdamore/tcell/v3"
)

// Cell represents a single character Cell on screen.
type Cell struct {
	Ch    rune // int32
	Style tcell.Style
	Size  int
	Width int // screen cell width
	Class CharClass
}

func (m *Cell) IsEmpty() bool {
	return m.Class == 0
}

func (c *Cell) Clear() {
	*c = Cell{}
}
