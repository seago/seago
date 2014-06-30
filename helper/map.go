package helper

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func MapToUrlValues(data map[string]interface{}) (url.Values, error) {
	var url_values = url.Values{}
	for k, v := range data {
		var s string
		switch v.(type) {
		case string:
			s = v.(string)
		case int:
			s = strconv.Itoa(v.(int))
		case int32:
			s = strconv.FormatInt(int64(v.(int32)), 10)
		case int64:
			s = strconv.FormatInt(v.(int64), 10)
		case float32:
			s = strconv.FormatFloat(float64(v.(float32)), 'f', 3, 64)
		case float64:
			s = strconv.FormatFloat(v.(float64), 'f', 3, 64)
		case []byte:
			s = string(v.([]byte))
		default:
			b, err := json.Marshal(v)
			if err == nil {
				s = string(b)
			}
			return nil, err
		}
		url_values.Add(k, s)
	}
	return url_values, nil
}
