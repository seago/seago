package seago

import (
	"encoding/json"
	"strconv"
	"text/template"
)

type Value string

/**
 * 获取string值
 */
func (v Value) String() string {
	return string(v)
}

/**
 * html escape
 */
func (v Value) HtmlEscape() string {
	return template.HTMLEscapeString(string(v))
}

/**
 * 获取int值
 */
func (v Value) Int() (i int) {
	i, _ = strconv.Atoi(string(v))
	return
}

/**
 * 获取int64值
 */
func (v Value) Int64() (i int64) {
	i, _ = strconv.ParseInt(string(v), 10, 64)
	return i
}

/**
 * 获取int32值
 */
func (v Value) Int32() int32 {
	i, _ := strconv.ParseInt(string(v), 10, 32)
	return int32(i)
}

/**
 * 获取float32值
 */
func (v Value) Float32() float32 {
	f, _ := strconv.ParseFloat(string(v), 32)
	return float32(f)
}

/**
 * 获取float64值
 */
func (v Value) Float64() float64 {
	f, _ := strconv.ParseFloat(string(v), 64)
	return f
}

/**
 * 获取bool值
 */
func (v Value) Bool() bool {
	b, _ := strconv.ParseBool(string(v))
	return b
}

/**
 * 解析json值
 */
func (v Value) JsonDecode(iv interface{}) error {
	return json.Unmarshal([]byte(string(v)), iv)
}
