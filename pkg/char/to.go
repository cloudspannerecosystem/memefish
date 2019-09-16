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
