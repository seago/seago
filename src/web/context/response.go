package context

import (
	"encoding/json"
	"mime"
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
	r.Header().Set(key, value)
}
func (r *Response) GetHeader(key string) string {
	return r.Header().Get(key)
}

func (r *Response) SetContentType(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	var content_type string

	if mtype := mime.TypeByExtension(ext); mtype != "" {
		content_type = mtype
	} else {
		content_type = "application/" + strings.TrimPrefix(ext, ".") + ";charset=utf-8"
	}
	r.SetHeader("Content-Type", content_type)
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
