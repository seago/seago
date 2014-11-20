package main

import (
	"github.com/seago/seago/web"
	"net/http"
)

func main() {
	m := web.New()
	m.Use(web.Logger())
	m.Use(web.Recovery())
	m.Use(web.Static("public", web.StaticOptions{
		Prefix:    "/public",
		IndexFile: "index.html",
	}))
	m.Before(func(w http.ResponseWriter, r *http.Request) bool {
		w.Header().Add("Server", "JDServer")
		return false
	})
	m.Get("/", func() string {
		return "index"
	})
	m.Get("/ping", func(c *web.Context) string {
		return "pong"
	})
	m.Run(8081)
}
