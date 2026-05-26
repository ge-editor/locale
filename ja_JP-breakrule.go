//go:build ja_JP

package locale

func IsBreakpoint(p2, p1, c Cell) bool {

	if p2.IsEmpty() || p1.IsEmpty() || c.IsEmpty() {
		return false
	}

	// 全て禁止文字
	if Is(p2, PROHIBITED) && Is(p1, PROHIBITED) && Is(c, PROHIBITED) {
		return true
	}

	// WORD なら1個の単語の構成とみなす
	if Is(p1, WORD) && Is(c, WORD) {
		return false
	}

	// 数字区切り文字の扱い e.g. 1,000,000
	isDecimalNumber :=
		Is(p2, NUMBER) && !Is(p2, WIDECHAR) &&
			Is(p1, DECIMAL_SEPARATOR) &&
			Is(c, NUMBER) && !Is(c, WIDECHAR)

	if Is(p1, PROHIBITED) &&
		!Is(c, PROHIBITED) &&
		!isDecimalNumber {
		return true
	}

	// Narrow, Wide 文字の違い
	if p1.Width != c.Width {
		return true
	}

	// カタカナ以外の文字種からカタカナが続いた場合
	if !Is(p1, KATAKANA) && Is(c, KATAKANA) {
		return true
	}

	return false
}
