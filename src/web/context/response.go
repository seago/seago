package context

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Response struct {
	http.ResponseWriter
}

func NewResponse(rw http.ResponseWriter) *Response {
	return &Response{rw}
}

func (r *Response) SetStatus(code int) {
	r.WriteHeader(code)
}

func (r *Response) SetHeader(key, value string) {
	r.Header().Set(strings.TrimSpace(key), strings.TrimSpace(value))
}

func (r *Response) GetHeader(key string) string {
	return r.Header().Get(strings.TrimSpace(key))
}

func (r *Response) WriteString(body string) (int, error) {
	return r.Write([]byte(body))
}

func (r *Response) WriteInternalServerError(b []byte) (int, error) {
	r.WriteHeader(http.StatusInternalServerError)
	return r.Write(b)
}

func (r *Response) WriteBadRequest(b []byte) (int, error) {
	r.WriteHeader(http.StatusBadRequest)
	return r.Write(b)
}

func (r *Response) JsonSuccess(v interface{}) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return r.WriteInternalServerError([]byte("json encoding error:" + err.Error()))
	}
	return r.Write(b)
}

func (r *Response) JsonError(v interface{}) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return r.WriteInternalServerError([]byte("json encoding error:" + err.Error()))
	}
	return r.WriteBadRequest(b)
}
