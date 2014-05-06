package seago

import (
	"github.com/seago/seago/server"
)

type Seago struct {
	Addr    string
	profile bool
	*server.Server
}

func NewSeago(addr string) *Seago {
	if addr == "" {
		addr = ":8080"
	}
	return &Seago{addr, false, server.NewServer()}
}

func (s *Seago) Profile() {
	s.profile = true
	s.Server.SetProfile(true)
}

func (s *Seago) Get(pattern string, hanlder interface{}) {
	s.Server.AddRouter(pattern, "GET", hanlder)
}

func (s *Seago) Put(pattern string, hanlder interface{}) {
	s.Server.AddRouter(pattern, "PUT", hanlder)
}

func (s *Seago) Post(pattern string, hanlder interface{}) {
	s.Server.AddRouter(pattern, "POST", hanlder)
}

func (s *Seago) Delete(pattern string, hanlder interface{}) {
	s.Server.AddRouter(pattern, "DELETE", hanlder)
}

func (s *Seago) Head(pattern string, hanlder interface{}) {
	s.Server.AddRouter(pattern, "HEAD", hanlder)
}

func (s *Seago) Options(pattern string, hanlder interface{}) {
	s.Server.AddRouter(pattern, "OPTIONS", hanlder)
}

func (s *Seago) Run() {
	s.Server.Run(s.Addr)
}
