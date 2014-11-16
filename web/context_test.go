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

package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/seago/seago/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Context(t *testing.T) {
	Convey("Do advanced encapsulation operations", t, func() {
		m := Classic()
		m.Use(Renderers(RenderOptions{
			Directory: "fixtures/basic",
		}, "fixtures/basic2"))

		Convey("Get request body", func() {
			m.Get("/body1", func(ctx *Context) {
				data, err := ioutil.ReadAll(ctx.Req.Body().ReadCloser())
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, "This is my request body")
			})
			m.Get("/body2", func(ctx *Context) {
				data, err := ctx.Req.Body().Bytes()
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, "This is my request body")
			})
			m.Get("/body3", func(ctx *Context) {
				data, err := ctx.Req.Body().String()
				So(err, ShouldBeNil)
				So(data, ShouldEqual, "This is my request body")
			})

			for i := 1; i <= 3; i++ {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/body"+utils.ToStr(i), nil)
				req.Body = ioutil.NopCloser(bytes.NewBufferString("This is my request body"))
				So(err, ShouldBeNil)
				m.ServeHTTP(resp, req)
			}
		})

		Convey("Get remote IP address", func() {
			m.Get("/remoteaddr", func(ctx *Context) string {
				return ctx.RemoteAddr()
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/remoteaddr", nil)
			req.RemoteAddr = "127.0.0.1:3333"
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "127.0.0.1")
		})

		Convey("Render HTML", func() {
			m.Get("/html", func(ctx *Context) {
				ctx.HTML(304, "hello", "Unknwon") // 304 for logger test.
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/html", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "<h1>Hello Unknwon</h1>")

			m.Get("/html2", func(ctx *Context) {
				ctx.Data["Name"] = "Unknwon"
				ctx.HTMLSet(200, "basic2", "hello2")
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/html2", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "<h1>Hello Unknwon</h1>")
		})

		Convey("Parse from and query", func() {
			m.Get("/query", func(ctx *Context) string {
				var buf bytes.Buffer
				buf.WriteString(ctx.Query("name") + " ")
				buf.WriteString(ctx.QueryEscape("name") + " ")
				buf.WriteString(utils.ToStr(ctx.QueryInt("int")) + " ")
				buf.WriteString(utils.ToStr(ctx.QueryInt64("int64")) + " ")
				return buf.String()
			})
			m.Get("/query2", func(ctx *Context) string {
				var buf bytes.Buffer
				buf.WriteString(strings.Join(ctx.QueryStrings("list"), ",") + " ")
				buf.WriteString(strings.Join(ctx.QueryStrings("404"), ",") + " ")
				return buf.String()
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/query?name=Unknwon&int=12&int64=123", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "Unknwon Unknwon 12 123 ")

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/query2?list=item1&list=item2", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "item1,item2  ")
		})

		Convey("URL parameter", func() {
			m.Get("/:name/:int/:int64", func(ctx *Context) string {
				var buf bytes.Buffer
				buf.WriteString(ctx.GetParam("name").String() + " ")
				buf.WriteString(ctx.GetParam("name").HtmlEscape() + " ")
				buf.WriteString(utils.ToStr(ctx.GetParam("int").Int()) + " ")
				buf.WriteString(utils.ToStr(ctx.GetParam("int64").Int64()) + " ")
				return buf.String()
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/user/1/13", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "user user 1 13 ")
		})

		Convey("Get file", func() {
			m.Get("/getfile", func(ctx *Context) {
				ctx.GetFile("hi")
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/getfile", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
		})

		Convey("Set and get cookie", func() {
			m.Get("/set", func(ctx *Context) {
				ctx.SetCookie("user", "Unknwon", 1)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/set", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Header().Get("Set-Cookie"), ShouldEqual, "user=Unknwon; Path=/; Max-Age=1")

			m.Get("/get", func(ctx *Context) string {
				ctx.GetCookie("404")
				return ctx.GetCookie("user")
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/get", nil)
			So(err, ShouldBeNil)
			req.Header.Set("Cookie", "user=Unknwon; Path=/; Max-Age=1")
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "Unknwon")
		})

		Convey("Set and get secure cookie", func() {
			m.SetDefaultCookieSecret("macaron")
			m.Get("/set", func(ctx *Context) {
				ctx.SetSecureCookie("user", "Unknwon", 1)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/set", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(strings.HasPrefix(resp.Header().Get("Set-Cookie"), "user=VW5rbndvbg==|"), ShouldBeTrue)

			m.Get("/get", func(ctx *Context) string {
				name, ok := ctx.GetSecureCookie("user")
				So(ok, ShouldBeTrue)
				return name
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/get", nil)
			So(err, ShouldBeNil)
			req.Header.Set("Cookie", "user=VW5rbndvbg==|1409244667158399419|6097781707f68d9940ba1ef0e78cc84aaeebc48f; Path=/; Max-Age=1")
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "Unknwon")
		})

		Convey("Serve files", func() {
			m.Get("/file", func(ctx *Context) {
				ctx.ServeFile("fixtures/custom_funcs/index.tmpl")
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/file", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "{{ myCustomFunc }}")

			m.Get("/file2", func(ctx *Context) {
				ctx.ServeFile("fixtures/custom_funcs/index.tmpl", "ok.tmpl")
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/file2", nil)
			So(err, ShouldBeNil)
			m.ServeHTTP(resp, req)
			So(resp.Body.String(), ShouldEqual, "{{ myCustomFunc }}")
		})
	})
}

func Test_Context_Redirect(t *testing.T) {
	Convey("Context with default redirect", t, func() {
		url, err := url.Parse("http://localhost/path/one")
		So(err, ShouldBeNil)
		resp := httptest.NewRecorder()
		req := http.Request{
			Method: "GET",
			URL:    url,
		}
		ctx := &Context{
			Req:     Request{&req},
			Resp:    NewResponseWriter(resp),
			Data:    make(map[string]interface{}),
			statics: make(map[string]*http.Dir),
		}
		ctx.Redirect("two")

		So(resp.Code, ShouldEqual, http.StatusFound)
		So(resp.HeaderMap["Location"][0], ShouldEqual, "/path/two")
	})

	Convey("Context with custom redirect", t, func() {
		url, err := url.Parse("http://localhost/path/one")
		So(err, ShouldBeNil)
		resp := httptest.NewRecorder()
		req := http.Request{
			Method: "GET",
			URL:    url,
		}
		ctx := &Context{
			Req:     Request{&req},
			Resp:    NewResponseWriter(resp),
			Data:    make(map[string]interface{}),
			statics: make(map[string]*http.Dir),
		}
		ctx.Redirect("two", 307)

		So(resp.Code, ShouldEqual, http.StatusTemporaryRedirect)
		So(resp.HeaderMap["Location"][0], ShouldEqual, "/path/two")
	})
}
