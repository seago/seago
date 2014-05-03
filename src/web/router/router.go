package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"web/context"
)

var NotFound = func(ctx *context.Context) {
	ctx.Response.SetStatus(http.StatusNotFound)
	ctx.Response.Write([]byte("Page Not Found"))
	return
}

type route struct {
	pattern     string
	regex       *regexp.Regexp
	handler     reflect.Value
	httpHandler http.Handler
	method      string
}

type Router struct {
	routes []*route
}

func NewRouter() *Router {
	return &Router{make([]*route, 0, 0)}
}

func (router *Router) AddRouter(pattern, method string, handler interface{}) error {
	r, err := newRoute(pattern, method, handler)
	if err != nil {
		return err
	}
	router.routes = append(router.routes, r)
	return nil
}

func (router *Router) Process(rw http.ResponseWriter, r *http.Request) {
	ctx := context.NewContext(rw, r)
	query := r.URL.Query()
	for k, _ := range query {
		ctx.SetParam(k, query.Get(k))
	}
	err := ctx.ParseForm(int64(2048000))
	if err != nil {
		log.Println(err)
	} else {
		if len(r.PostForm) > 0 {
			query = r.PostForm
			for k, _ := range query {
				ctx.SetParam(k, query.Get(k))
			}
		}
		if r.MultipartForm != nil {
			query = r.MultipartForm.Value
			for k, v := range query {
				ctx.SetParam(k, v[0])
			}
			for k, v := range r.MultipartForm.File {
				ctx.File[k] = v[0]
			}
		}
	}
	for k, _ := range query {
		ctx.SetParam(k, query.Get(k))
	}

	route := router.process(ctx)
	if route != nil {
		route.httpHandler.ServeHTTP(rw, r)
	}
}

func (router *Router) process(ctx *context.Context) (unused *route) {
	path := ctx.Request.URL.Path

	for _, route := range router.routes {

		if !route.regex.MatchString(path) {
			continue
		}
		ok, params := route.match(ctx.Request.Method, path)
		if !ok {
			continue
		}
		if route.httpHandler != nil {
			unused = route
			return
		}

		var args []reflect.Value
		handlerType := route.handler.Type()
		if requireContext(handlerType) {
			args = append(args, reflect.ValueOf(ctx))
		}
		if len(params) > 0 {
			for _, v := range params {
				args = append(args, reflect.ValueOf(v))
			}
		}

		ret := route.handler.Call(args)
		if len(ret) == 0 {
			return
		}

		sval := ret[0]
		var content []byte
		if sval.Kind() == reflect.String {
			content = []byte(sval.String())
		} else if sval.Kind() == reflect.Slice && sval.Type().Elem().Kind() == reflect.Uint8 {
			content = sval.Interface().([]byte)
		}
		ctx.Response.SetHeader("Content-Length", strconv.Itoa(len(content)))
		_, err := ctx.Response.Write(content)
		if err != nil {
			log.Println("Error response write")
		}
		return
	}
	//TODO:Page not found
	NotFound(ctx)
	return
}

func newRoute(pattern, method string, handler interface{}) (*route, error) {

	route := &route{pattern: pattern, method: strings.ToUpper(method)}
	r := regexp.MustCompile(`:[^/#?()\.\\]+`)
	pattern = r.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	r2 := regexp.MustCompile(`\*\*`)
	var index int
	pattern = r2.ReplaceAllStringFunc(pattern, func(m string) string {
		index++
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	})
	pattern += `\/?`
	route.regex = regexp.MustCompile(pattern)
	switch handler.(type) {
	case http.Handler:
		route.httpHandler = handler.(http.Handler)
	case reflect.Value:
		fv := handler.(reflect.Value)
		route.handler = fv
	default:
		if !checkHandler(handler) {
			return nil, errors.New("handler is not func")
		}
		fv := reflect.ValueOf(handler)
		route.handler = fv
	}
	return route, nil
}

func (r *route) matchMethod(method string) bool {
	return r.method == "*" || method == r.method || (method == "HEAD" && r.method == "GET")
}

func (r *route) match(method, path string) (bool, []string) {
	if !r.matchMethod(method) {
		return false, nil
	}
	matchs := r.regex.FindStringSubmatch(path)
	if len(matchs) > 0 && matchs[0] == path {
		return true, matchs[1:]
	}
	return false, nil
}

func requireContext(handlerType reflect.Type) bool {
	if handlerType.NumIn() == 0 {
		return false
	}

	arg0 := handlerType.In(0)

	if arg0.Kind() != reflect.Ptr {
		return false
	}

	if arg0.Elem() == reflect.TypeOf(context.Context{}) {
		return true
	}
	return false
}

func checkHandler(handler interface{}) bool {
	return reflect.ValueOf(handler).Kind() == reflect.Func
}
