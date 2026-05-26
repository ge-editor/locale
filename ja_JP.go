//go:build ja_JP

package locale

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type ja_JP struct{}

func New() Locale {
	return &ja_JP{}
}

func (l *ja_JP) LocaleName() string {
	return "ja_JP"
}

func (l *ja_JP) DisplayName() string {
	return "Japanese"
}

// Character Classification

func (l *ja_JP) GetCharClass(ch rune) (cc CharClass) {
	prohibitedCharacters := "\t !&)*+,-./:;<=>?@]|}~　、。」）！？をん"
	if strings.ContainsRune(prohibitedCharacters, ch) {
		cc |= PROHIBITED
	}

	switch {
	// 1. 個別の特殊文字・制御文字（最優先）
	case ch == ',' || ch == '.':
		cc |= DECIMAL_SEPARATOR
	case ch == '\t':
		cc |= CONTROLCODE | TAB
	case ch == ' ':
		cc |= SPACE
	case ch == '　':
		cc |= SPACE | WIDECHAR
	case ch == 127: // DEL character
		cc |= CONTROLCODE | DEL
	case ch == '\n':
		cc |= CONTROLCODE | NEWLINE
	case ch == 26: // EOF (typically Ctrl+Z)
		cc |= CONTROLCODE | EOF
	case ch < 32:
		cc |= CONTROLCODE

	// 2. 日本語固有の文字（ひらがな・カタカナ・漢字）
	// ※ IsLetter より上に配置することで、日本語として正しく分類する
	case unicode.Is(unicode.Hiragana, ch):
		cc |= HIRAGANA | WIDECHAR
		// gelog.Info("HIRAGANA", string(ch))
	case unicode.Is(unicode.Katakana, ch):
		cc |= KATAKANA | WIDECHAR
		// gelog.Info("KATAKANA", string(ch))
	case ch >= 0x4E00 && ch <= 0x9FFF: // CJK統合漢字の範囲
		cc |= WIDECHAR

	// 3. 全角英数字
	// ※ 汎用の IsDigit/IsLetter より上に配置して WIDECHAR フラグを確実に立てる
	case ch >= '０' && ch <= '９': // 全角数字
		cc |= NUMBER | WIDECHAR
	case ch >= 'Ａ' && ch <= 'Ｚ': // 全角アルファベット（大文字）
		cc |= ALPHABET | WIDECHAR | UPPERCASE
	case ch >= 'ａ' && ch <= 'ｚ': // 全角アルファベット（小文字）
		cc |= ALPHABET | WIDECHAR

	// 4. 汎用的な文字・数字判定（半角英数字などがここにヒットする）
	case unicode.IsDigit(ch):
		cc |= NUMBER | WORD
	case unicode.IsLetter(ch):
		cc |= ALPHABET | WORD
		if unicode.IsUpper(ch) {
			cc |= UPPERCASE
		}
	case ch == '_':
		cc |= SYMBOL | WORD
	// case unicode.IsSymbol(ch) || unicode.IsPunct(ch):
	case unicode.IsSymbol(ch):
		cc |= SYMBOL

	default:
		cc |= OTHER
	}

	return cc
}

// Width

func (l *ja_JP) RuneWidth(r rune) int {
	return eastAsianWidth(r, true)
}

// Bullets

func (l *ja_JP) Bullets() []Bullet {
	return []Bullet{
		{Marker: []byte("- "), Width: 2},
		{Marker: []byte("* "), Width: 2},
		{Marker: []byte("・"), Width: 2},
		{Marker: []byte("●"), Width: 2},
		{Marker: []byte("■"), Width: 2},
	}
}

// Break Policy

type jaBreakPolicy struct {
	rules []BreakRule
}

func (l *ja_JP) BreakPolicy() BreakPolicy {
	return &jaBreakPolicy{
		rules: []BreakRule{
			ruleClassChange,
		},
	}
}

func (p *jaBreakPolicy) IsBreakPoint(p2, p1, c rune) bool {
	for _, r := range p.rules {
		if r(p2, p1, c) {
			return true
		}
	}
	return false
}

func ruleClassChange(p2, p1, c rune) bool {
	// simple rule: class change creates boundary
	return p1 != c
}

// Word Policy

type jaWordPolicy struct {
	locale *ja_JP
}

func (l *ja_JP) WordPolicy() WordPolicy {
	return &jaWordPolicy{locale: l}
}

func (p *jaWordPolicy) Boundary(text []byte, cursor int) (int, int, bool) {
	if cursor < 0 || cursor >= len(text) {
		return 0, 0, false
	}

	r, size := utf8.DecodeRune(text[cursor:])
	if r == utf8.RuneError {
		return 0, 0, false
	}

	class := p.locale.GetCharClass(r)
	if class&SPACE != 0 || class&SYMBOL != 0 {
		return 0, 0, false
	}

	start := cursor
	end := cursor + size

	// backward
	for start > 0 {
		r2, s2 := utf8.DecodeLastRune(text[:start])
		if p.locale.GetCharClass(r2) != class {
			break
		}
		start -= s2
	}

	// forward
	for end < len(text) {
		r2, s2 := utf8.DecodeRune(text[end:])
		if p.locale.GetCharClass(r2) != class {
			break
		}
		end += s2
	}

	return start, end, true
}
