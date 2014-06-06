package utils

func Iter(i int) []struct{} {
	return make([]struct{}, i)
}

func MapToList(m map[interface{}]interface{}) []interface{} {
	l := make([]interface{}, 0)
	for _, v := range m {
		l = append(l, v)
	}
	return l
}
