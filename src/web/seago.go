package web

import (
	"web/server"
)

type Seago struct {
	Addr    string
	profile bool
	server  *server.Server
}

func NewSeago(addr string) *Seago {
	if addr == "" {
		addr = ":8080"
	}
	return &Seago{addr, false, server.NewServer()}
}

func (s *Seago) Profile() {
	s.profile = true
	s.server.SetProfile(true)
}

func (s *Seago) Get(pattern string, hanlder interface{}) {
	s.server.AddRouter(pattern, "GET", hanlder)
}

func (s *Seago) Put(pattern string, hanlder interface{}) {
	s.server.AddRouter(pattern, "PUT", hanlder)
}

func (s *Seago) Post(pattern string, hanlder interface{}) {
	s.server.AddRouter(pattern, "POST", hanlder)
}

func (s *Seago) Delete(pattern string, hanlder interface{}) {
	s.server.AddRouter(pattern, "DELETE", hanlder)
}

func (s *Seago) Head(pattern string, hanlder interface{}) {
	s.server.AddRouter(pattern, "HEAD", hanlder)
}

func (s *Seago) Options(pattern string, hanlder interface{}) {
	s.server.AddRouter(pattern, "OPTIONS", hanlder)
}

func (s *Seago) Run() {
	s.server.Run(s.Addr)
}
