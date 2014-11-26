package main

import (
	"github.com/seago/seago"
	"net/http"
)

func main() {
	s := seago.New()
	s.Use(seago.Logger())
	s.Use(seago.Recovery())
	s.Use(seago.Static("public", seago.StaticOptions{
		Prefix:    "/public",
		IndexFile: "index.html",
	}))
	s.Before(func(w http.ResponseWriter, r *http.Request) bool {
		w.Header().Add("Server", "JDServer")
		return false
	})
	s.Get("/ping/user_:id", func() string {
		return "index"
	})

	s.Get("/ping/test/:id", func(c *seago.Context) string {
		return "pong"
	})
	s.Run(8081)
}
