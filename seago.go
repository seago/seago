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

// Package Seago is a high productive and modular design web framework in Go.
package seago

import (
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/seago/com"
	"github.com/seago/seago/inject"

	"gopkg.in/ini.v1"
)

const _VERSION = "0.0.2.0722"

func Version() string {
	return _VERSION
}

// Handler can be any callable function.
// Seago attempts to inject services into the handler's argument list,
// and panics if an argument could not be fullfilled via dependency injection.
type Handler interface{}

// validateHandler makes sure a handler is a callable function,
// and panics if it is not.
func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("Web handler must be a callable function")
	}
}

// Seago represents the top level web application.
// inject.Injector methods can be invoked to map services on a global level.
type Seago struct {
	inject.Injector
	befores  []BeforeHandler
	handlers []Handler
	action   Handler

	urlPrefix string // For suburl support.
	*Router

	logger *log.Logger
}

// NewWithLogger creates a bare bones Seago instance.
// Use this method if you want to have full control over the middleware that is used.
// You can specify logger output writer with this function.
func NewWithLogger(out io.Writer) *Seago {
	m := &Seago{
		Injector: inject.New(),
		action:   func() {},
		Router:   NewRouter(),
		logger:   log.New(out, "[SEAGO] ", 0),
	}
	m.Router.m = m
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	m.notFound = func(resp http.ResponseWriter, req *http.Request) {
		c := m.createContext(resp, req)
		c.handlers = append(c.handlers, http.NotFound)
		c.run()
	}
	return m
}

// New creates a bare bones Seago instance.
// Use this method if you want to have full control over the middleware that is used.
func New() *Seago {
	return NewWithLogger(os.Stdout)
}

// Classic creates a classic Macaron with some basic default middleware:
// Web.Logger, Web.Recovery and Web.Static.
func Classic() *Seago {
	m := New()
	m.Use(Logger())
	m.Use(Recovery())
	m.Use(Static("public"))
	return m
}

// Handlers sets the entire middleware stack with the given Handlers.
// This will clear any current middleware handlers,
// and panics if any of the handlers is not a callable function
func (m *Seago) Handlers(handlers ...Handler) {
	m.handlers = make([]Handler, 0)
	for _, handler := range handlers {
		m.Use(handler)
	}
}

// Action sets the handler that will be called after all the middleware has been invoked.
// This is set to macaron.Router in a macaron.Classic().
func (m *Seago) Action(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

// BeforeHandler represents a handler executes at beginning of every request.
// Macaron stops future process when it returns true.
type BeforeHandler func(rw http.ResponseWriter, req *http.Request) bool

func (m *Seago) Before(handler BeforeHandler) {
	m.befores = append(m.befores, handler)
}

// Use adds a middleware Handler to the stack,
// and panics if the handler is not a callable func.
// Middleware Handlers are invoked in the order that they are added.
func (m *Seago) Use(handler Handler) {
	validateHandler(handler)
	m.handlers = append(m.handlers, handler)
}

func (m *Seago) createContext(rw http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		Injector: inject.New(),
		handlers: m.handlers,
		action:   m.action,
		index:    0,
		Router:   m.Router,
		Request:  Request{req},
		Response: NewResponseWriter(rw),
		Data:     make(map[string]interface{}),
		statics:  make(map[string]*http.Dir),
	}
	c.SetParent(m)
	c.Map(c)
	c.MapTo(c.Response, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

// ServeHTTP is the HTTP Entry point for a Macaron instance.
// Useful if you want to control your own HTTP server.
// Be aware that none of middleware will run without registering any router.
func (m *Seago) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.TrimPrefix(req.URL.Path, m.urlPrefix)
	for _, h := range m.befores {
		if h(rw, req) {
			return
		}
	}
	m.Router.ServeHTTP(rw, req)
}

func GetDefaultListenInfo() (string, int) {
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	port := com.StrTo(os.Getenv("PORT")).MustInt()
	if port == 0 {
		port = 4000
	}
	return host, port
}

// Run the http server. Listening on os.GetEnv("PORT") or 4000 by default.
func (m *Seago) Run(args ...interface{}) {
	host, port := GetDefaultListenInfo()
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case string:
			host = arg
		case int:
			port = arg
		}
	} else if len(args) >= 2 {
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}

	addr := host + ":" + com.ToStr(port)
	logger := m.Injector.GetVal(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("listening on %s (%s)\n", addr, Env)
	logger.Fatalln(http.ListenAndServe(addr, m))
}

// SetURLPrefix sets URL prefix of router layer, so that it support suburl.
func (m *Seago) SetURLPrefix(prefix string) {
	m.urlPrefix = prefix
}

// ____   ____            .__      ___.   .__
// \   \ /   /____ _______|__|____ \_ |__ |  |   ____   ______
//  \   Y   /\__  \\_  __ \  \__  \ | __ \|  | _/ __ \ /  ___/
//   \     /  / __ \|  | \/  |/ __ \| \_\ \  |_\  ___/ \___ \
//    \___/  (____  /__|  |__(____  /___  /____/\___  >____  >
//                \/              \/    \/          \/     \/

const (
	DEV  string = "development"
	PROD string = "production"
	TEST string = "test"
)

var (
	// Env is the environment that Macaron is executing in.
	// The SEAGO_ENV is read on initialization to set this variable.
	Env = DEV

	// Path of work directory.
	Root string

	// Flash applies to current request.
	FlashNow bool
	// Configuration convention object.
	cfg *ini.File
)

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	setENV(os.Getenv("SEAGO_ENV"))
	var err error
	Root, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

// SetConfig sets data sources for configuration.
func SetConfig(source interface{}, others ...interface{}) (_ *ini.File, err error) {
	cfg, err = ini.Load(source, others...)
	return Config(), err
}

// Config returns configuration convention object.
// It returns an empty object if there is no one available.
func Config() *ini.File {
	if cfg == nil {
		return ini.Empty()
	}
	return cfg
}
