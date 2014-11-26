// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package seago

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/seago/seago/inject"
)

// Locale reprents a localization interface.
type Locale interface {
	Language() string
	Tr(string, ...interface{}) string
}

// Body is the request's body.
type RequestBody struct {
	reader io.ReadCloser
}

func (rb *RequestBody) Bytes() ([]byte, error) {
	return ioutil.ReadAll(rb.reader)
}

func (rb *RequestBody) String() (string, error) {
	data, err := rb.Bytes()
	return string(data), err
}

func (rb *RequestBody) ReadCloser() io.ReadCloser {
	return rb.reader
}

// A Request represents an HTTP request received by a server or to be sent by a client.
type Request struct {
	*http.Request
}

func (r *Request) Body() *RequestBody {
	return &RequestBody{r.Request.Body}
}

// Context represents the runtime context of current request of Seago instance.
// It is the integration of most frequently used middlewares and helper methods.
type Context struct {
	inject.Injector
	handlers []Handler
	action   Handler
	index    int

	*Router
	Request  Request
	Response ResponseWriter
	params   Params

	File   *multipart.File
	Render // Not nil only if you use macaran.Render middleware.
	Locale
	Data    map[string]interface{}
	statics map[string]*http.Dir // Registered static paths.
}

func (c *Context) handler() Handler {
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("invalid index for context handler")
}

func (c *Context) Next() {
	c.index += 1
	c.run()
}

func (c *Context) Written() bool {
	return c.Response.Written()
}

func (c *Context) run() {
	for c.index <= len(c.handlers) {
		vals, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index += 1

		// if the handler returned something, write it to the http response
		if len(vals) > 0 {
			ev := c.GetVal(reflect.TypeOf(ReturnHandler(nil)))
			handleReturn := ev.Interface().(ReturnHandler)
			handleReturn(c, vals)
		}

		if c.Written() {
			return
		}
	}
}

// RemoteAddr returns more real IP address.
func (ctx *Context) RemoteAddr() string {
	addr := ctx.Request.Header.Get("X-Real-IP")
	if len(addr) == 0 {
		addr = ctx.Request.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = ctx.Request.RemoteAddr
			if i := strings.LastIndex(addr, ":"); i > -1 {
				addr = addr[:i]
			}
		}
	}
	return addr
}

func (ctx *Context) renderHTML(status int, setName, tplName string, data ...interface{}) {
	if ctx.Render == nil {
		panic("renderer middleware hasn't been registered")
	}
	if len(data) == 0 {
		ctx.Render.HTMLSet(status, setName, tplName, ctx.Data)
	} else {
		ctx.Render.HTMLSet(status, setName, tplName, data[0])
		if len(data) > 1 {
			ctx.Render.HTMLSet(status, setName, tplName, data[0], data[1].(HTMLOptions))
		}
	}
}

// HTML calls Render.HTML but allows less arguments.
func (ctx *Context) HTML(status int, name string, data ...interface{}) {
	ctx.renderHTML(status, _DEFAULT_TPL_SET_NAME, name, data...)
}

// HTML calls Render.HTMLSet but allows less arguments.
func (ctx *Context) HTMLSet(status int, setName, tplName string, data ...interface{}) {
	ctx.renderHTML(status, setName, tplName, data...)
}

f//URL跳转
func (ctx *Context) Redirect(location string, status ...int) {
	code := http.StatusFound
	if len(status) == 1 {
		code = status[0]
	}

	http.Redirect(ctx.Response, ctx.Request.Request, location, code)
}

//获取http 请求中 GET POST PUT DELETE等方法的参数以及路由规则中的参数
func (ctx *Context) GetParam(name string) Value {
	return Value(ctx.params[name])
}

// GetFile returns information about user upload file by given form field name.
func (ctx *Context) GetFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.Request.FormFile(name)
}

// SetCookie sets given cookie value to response header.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = value

	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			cookie.MaxAge = v
		case int64:
			cookie.MaxAge = int(v)
		case int32:
			cookie.MaxAge = int(v)
		}
	}

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	} else {
		cookie.Path = "/"
	}

	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	// default empty
	if len(others) > 3 {
		switch v := others[3].(type) {
		case bool:
			cookie.Secure = v
		default:
			if others[3] != nil {
				cookie.Secure = true
			}
		}
	}

	// default false. for session cookie default true
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			cookie.HttpOnly = true
		}
	}

	ctx.Response.Header().Add("Set-Cookie", cookie.String())
}

// GetCookie returns given cookie value from request header.
func (ctx *Context) GetCookie(name string) string {
	cookie, err := ctx.Request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

var defaultCookieSecret string

// SetDefaultCookieSecret sets global default secure cookie secret.
func (m *Seago) SetDefaultCookieSecret(secret string) {
	defaultCookieSecret = secret
}

// SetSecureCookie sets given cookie value to response header with default secret string.
func (ctx *Context) SetSecureCookie(name, value string, others ...interface{}) {
	ctx.SetSuperSecureCookie(defaultCookieSecret, name, value, others...)
}

// GetSecureCookie returns given cookie value from request header with default secret string.
func (ctx *Context) GetSecureCookie(key string) (string, bool) {
	return ctx.GetSuperSecureCookie(defaultCookieSecret, key)
}

// SetSuperSecureCookie sets given cookie value to response header with secret string.
func (ctx *Context) SetSuperSecureCookie(Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.SetCookie(name, cookie, others...)
}

// GetSuperSecureCookie returns given cookie value from request header with secret string.
func (ctx *Context) GetSuperSecureCookie(Secret, key string) (string, bool) {
	val := ctx.GetCookie(key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// ServeFile serves given file to response.
func (ctx *Context) ServeFile(file string, names ...string) {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		name = path.Base(file)
	}
	ctx.Response.Header().Set("Content-Description", "File Transfer")
	ctx.Response.Header().Set("Content-Type", "application/octet-stream")
	ctx.Response.Header().Set("Content-Disposition", "attachment; filename="+name)
	ctx.Response.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Response.Header().Set("Expires", "0")
	ctx.Response.Header().Set("Cache-Control", "must-revalidate")
	ctx.Response.Header().Set("Pragma", "public")
	http.ServeFile(ctx.Response, ctx.Request.Request, file)
}

// ChangeStaticPath changes static path from old to new one.
func (ctx *Context) ChangeStaticPath(oldPath, newPath string) {
	if !filepath.IsAbs(oldPath) {
		oldPath = filepath.Join(Root, oldPath)
	}
	if _, ok := ctx.statics[oldPath]; ok {
		dir := ctx.statics[oldPath]
		delete(ctx.statics, oldPath)

		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(Root, newPath)
		}
		*dir = http.Dir(newPath)
		ctx.statics[newPath] = dir
	}
}
