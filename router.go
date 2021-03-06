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
	"net/http"
	"strings"
	"sync"

	"github.com/seago/com"
)

var (
	// Known HTTP methods.
	_HTTP_METHODS = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"OPTIONS": true,
		"HEAD":    true,
	}
)

// routeMap represents a thread-safe map for route tree.
type routeMap struct {
	lock   sync.RWMutex
	routes map[string]map[string]bool
}

// NewRouteMap initializes and returns a new routeMap.
func NewRouteMap() *routeMap {
	rm := &routeMap{
		routes: make(map[string]map[string]bool),
	}
	for m := range _HTTP_METHODS {
		rm.routes[m] = make(map[string]bool)
	}
	return rm
}

// isExist returns true if a route has been registered.
func (rm *routeMap) isExist(method, pattern string) bool {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	return rm.routes[method][pattern]
}

// add adds new route to route tree map.
func (rm *routeMap) add(method, pattern string) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	rm.routes[method][pattern] = true
}

type group struct {
	pattern  string
	handlers []Handler
}

// Router represents a Seago router layer.
type Router struct {
	m       *Seago
	routers map[string]*Tree
	*routeMap

	groups   []group
	notFound http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routers:  make(map[string]*Tree),
		routeMap: NewRouteMap(),
	}
}

type Params map[string]string

// Handle is a function that can be registered to a route to handle HTTP requests.
// Like http.HandlerFunc, but has a third parameter for the values of wildcards (variables).
type Handle func(http.ResponseWriter, *http.Request, Params)

// handle adds new route to the router tree.
func (r *Router) handle(method, pattern string, handle Handle) {
	method = strings.ToUpper(method)

	// Prevent duplicate routes.
	if r.isExist(method, pattern) {
		return
	}

	// Validate HTTP methods.
	if !_HTTP_METHODS[method] && method != "*" {
		panic("unknown HTTP method: " + method)
	}

	// Generate methods need register.
	methods := make(map[string]bool)
	if method == "*" {
		for m := range _HTTP_METHODS {
			methods[m] = true
		}
	} else {
		methods[method] = true
	}

	// Add to router tree.
	for m := range methods {
		if t, ok := r.routers[m]; ok {
			t.AddRouter(pattern, handle)
		} else {
			t := NewTree()
			t.AddRouter(pattern, handle)
			r.routers[m] = t
		}
		r.add(m, pattern)
	}
}

// Handle registers a new request handle with the given pattern, method and handlers.
func (r *Router) Handle(method string, pattern string, handlers []Handler) {
	if len(r.groups) > 0 {
		groupPattern := ""
		h := make([]Handler, 0)
		for _, g := range r.groups {
			groupPattern += g.pattern
			h = append(h, g.handlers...)
		}

		pattern = groupPattern + pattern
		h = append(h, handlers...)
		handlers = h
	}
	// verify handlers by cnphpbb at 20140803 23:51
	for _, handler := range handlers {
		validateHandler(handler)
	}

	r.handle(method, pattern, func(resp http.ResponseWriter, req *http.Request, params Params) {
		c := r.m.createContext(resp, req)
		if len(params) == 0 {
			c.params = make(map[string]string)
		} else {
			c.params = params
		}

		//处理get请求的参数给context's params
		query := req.URL.Query()
		for k := range query {
			c.params[k] = query.Get(k)
		}
		//处理post请求的参数赋值给context's params
		if method == "POST" || method == "PUT" || method == "PATCH" {
			if req.Form == nil {
				if strings.Contains(req.Header.Get("Content-Type"), "multipart/form-data") {
					req.ParseMultipartForm(32 << 20)
				} else {
					req.ParseForm()
				}
				if len(req.PostForm) > 0 {
					query = req.PostForm
					for k, _ := range query {
						c.params[k] = query.Get(k)
					}
				}
				if len(req.Form) > 0 {
					query = req.Form
					for k, _ := range query {
						c.params[k] = query.Get(k)
					}
				}
			}
		}
		c.handlers = append(r.m.handlers, handlers...)
		c.run()
	})
}

func (r *Router) Group(pattern string, fn func(), h ...Handler) {
	r.groups = append(r.groups, group{pattern, h})
	fn()
	r.groups = r.groups[:len(r.groups)-1]
}

// Get is a shortcut for r.Handle("GET", pattern, handlers)
func (r *Router) Get(pattern string, h ...Handler) {
	r.Handle("GET", pattern, h)
}

// Patch is a shortcut for r.Handle("PATCH", pattern, handlers)
func (r *Router) Patch(pattern string, h ...Handler) {
	r.Handle("PATCH", pattern, h)
}

// Post is a shortcut for r.Handle("POST", pattern, handlers)
func (r *Router) Post(pattern string, h ...Handler) {
	r.Handle("POST", pattern, h)
}

// Put is a shortcut for r.Handle("PUT", pattern, handlers)
func (r *Router) Put(pattern string, h ...Handler) {
	r.Handle("PUT", pattern, h)
}

// Delete is a shortcut for r.Handle("DELETE", pattern, handlers)
func (r *Router) Delete(pattern string, h ...Handler) {
	r.Handle("DELETE", pattern, h)
}

// Options is a shortcut for r.Handle("OPTIONS", pattern, handlers)
func (r *Router) Options(pattern string, h ...Handler) {
	r.Handle("OPTIONS", pattern, h)
}

// Head is a shortcut for r.Handle("HEAD", pattern, handlers)
func (r *Router) Head(pattern string, h ...Handler) {
	r.Handle("HEAD", pattern, h)
}

// Any is a shortcut for r.Handle("*", pattern, handlers)
func (r *Router) Any(pattern string, h ...Handler) {
	r.Handle("*", pattern, h)
}

// Route is a shortcut for same handlers but different HTTP methods.
//
// Example:
// 		m.Route("/", "GET,POST", h)
func (r *Router) Route(pattern, methods string, h ...Handler) {
	for _, m := range strings.Split(methods, ",") {
		r.Handle(strings.TrimSpace(m), pattern, h)
	}
}

// Combo returns a combo router.
func (r *Router) Combo(pattern string) *ComboRouter {
	return &ComboRouter{r, pattern, map[string]bool{}}
}

// Configurable http.HandlerFunc which is called when no matching route is
// found. If it is not set, http.NotFound is used.
// Be sure to set 404 response code in your handler.
func (r *Router) NotFound(handlers ...Handler) {
	r.notFound = func(rw http.ResponseWriter, req *http.Request) {
		c := r.m.createContext(rw, req)
		c.handlers = append(r.m.handlers, handlers...)
		c.run()
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if t, ok := r.routers[req.Method]; ok {
		h, p := t.Match(req.URL.Path)
		if h != nil {
			if splat, ok := p[":splat"]; ok {
				p["*"] = p[":splat"] // Better name.
				splatlist := strings.Split(splat, "/")
				for k, v := range splatlist {
					p[com.ToStr(k)] = v
				}
			}
			h(rw, req, p)
			return
		}
	}

	r.notFound(rw, req)
}

// ComboRouter represents a combo router.
type ComboRouter struct {
	router  *Router
	pattern string
	methods map[string]bool // Registered methods.
}

func (cr *ComboRouter) checkMethod(name string) {
	if cr.methods[name] {
		panic("method '" + name + "' has already been registered")
	}
	cr.methods[name] = true
}

func (cr *ComboRouter) Get(h ...Handler) *ComboRouter {
	cr.checkMethod("GET")
	cr.router.Get(cr.pattern, h...)
	return cr
}

func (cr *ComboRouter) Patch(h ...Handler) *ComboRouter {
	cr.checkMethod("PATCH")
	cr.router.Patch(cr.pattern, h...)
	return cr
}

func (cr *ComboRouter) Post(h ...Handler) *ComboRouter {
	cr.checkMethod("POST")
	cr.router.Post(cr.pattern, h...)
	return cr
}

func (cr *ComboRouter) Put(h ...Handler) *ComboRouter {
	cr.checkMethod("PUT")
	cr.router.Put(cr.pattern, h...)
	return cr
}

func (cr *ComboRouter) Delete(h ...Handler) *ComboRouter {
	cr.checkMethod("DELETE")
	cr.router.Delete(cr.pattern, h...)
	return cr
}

func (cr *ComboRouter) Options(h ...Handler) *ComboRouter {
	cr.checkMethod("OPTIONS")
	cr.router.Options(cr.pattern, h...)
	return cr
}

func (cr *ComboRouter) Head(h ...Handler) *ComboRouter {
	cr.checkMethod("HEAD")
	cr.router.Head(cr.pattern, h...)
	return cr
}
