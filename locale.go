package locale

import (
	"unicode"

	"golang.org/x/text/width"
)

// Locale defines cultural text behavior.
// Internal text must be UTF-8.

type Locale interface {
	LocaleName() string
	DisplayName() string

	GetCharClass(r rune) CharClass
	RuneWidth(r rune) int

	Bullets() []Bullet
	BreakPolicy() BreakPolicy
	WordPolicy() WordPolicy
}

// Character Class

type CharClass int

const (
	OTHER CharClass = 1 << iota
	CONTROLCODE
	TAB
	NEWLINE
	DEL
	EOF
	WORD // Treat it as a set of words
	NUMBER
	ALPHABET
	PROHIBITED
	DECIMAL_SEPARATOR
	WIDECHAR // 全角文字
	UPPERCASE
	SYMBOL
	SPACE
	HIRAGANA
	KATAKANA
)

// Bullet

type Bullet struct {
	Marker []byte
	Width  int
}

// Break Policy

type BreakRule func(p2, p1, c rune) bool

type BreakPolicy interface {
	IsBreakPoint(p2, p1, c rune) bool
}

// Word Policy

type WordPolicy interface {
	Boundary(text []byte, cursor int) (start, end int, ok bool)
}

// East Asian Width

func eastAsianWidth(r rune, ambiguousWide bool) int {
	if unicode.Is(unicode.Mn, r) {
		return 0
	}

	k := width.LookupRune(r).Kind()

	switch k {
	case width.EastAsianWide,
		width.EastAsianFullwidth:
		return 2

	case width.EastAsianAmbiguous:
		if ambiguousWide {
			return 2
		}
		return 1

	case width.EastAsianHalfwidth,
		width.EastAsianNarrow:
		return 1
	}

	return 1
}

func Is(c Cell, flag CharClass) bool {
	return c.Class&flag != 0
}
