package main

import (
	"github.com/seago/seago/web"
)

func main() {
	m := web.New()
	m.Use(web.Logger())
	m.Use(web.Recovery())
	m.Use(web.Static("public", web.StaticOptions{
		Prefix:    "/public",
		IndexFile: "index.html",
	}))
	m.Before(func(c *web.Context) {
		c.Header().Add("Server", "JDServer")
	})
	m.Get("/", func() string {
		return "index"
	})
	m.Get("/ping", func(c *web.Context) string {
		return "pong"
	})
	m.Run(8081)
}
