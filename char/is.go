package char

func IsPrint(b byte) bool {
	return 0x20 <= b && b <= 0x7E // ' ' to '~'
}

func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsHexDigit(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func IsOctalDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

func IsIdentStart(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func IsIdentPart(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}
