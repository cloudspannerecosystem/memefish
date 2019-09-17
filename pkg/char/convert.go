package char

func ToUpper(s string) string {
	var bs []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			if bs == nil {
				bs = []byte(s[:i])
			}
			bs = append(bs, 'A'+(c-'a'))
		} else if bs != nil {
			bs = append(bs, c)
		}
	}

	if bs == nil {
		return s
	}
	return string(bs)
}

func EqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}

	for i := 0; i < len(s); i++ {
		c, d := s[i], t[i]
		if 'a' <= c && c <= 'z' {
			c = 'A' + (c - 'a')
		}
		if 'a' <= d && d <= 'z' {
			d = 'A' + (d - 'a')
		}
		if c != d {
			return false
		}
	}

	return true
}
