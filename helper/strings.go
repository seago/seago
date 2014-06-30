package helper

//字符串切割，支持中文
func SubStr(str string, start, length int) (substr string) {
	rs := []rune(str)
	rune_length := len(rs)
	if start < 0 {
		start = 0
	}

	if start >= rune_length {
		return ""
	}
	end := start + length
	if end > rune_length {
		end = rune_length
	}
	return string(rs[start:end])
}

func StrLen(str string) int {
	rs := []rune(str)
	return len(rs)
}
