package server

import (
	//"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	//"web/context"
	"log"
	"os"
	"web/router"
)

type Server struct {
	router   *router.Router
	Logger   *log.Logger
	l        net.Listener
	profiler bool
}

func NewServer() *Server {
	return &Server{
		router: router.NewRouter(),
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (s *Server) SetProfile(is bool) {
	s.profiler = is
}

func (s *Server) AddRouter(pattern, method string, handler interface{}) {
	err := s.router.AddRouter(pattern, method, handler)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.Process(rw, r)
}

func (s *Server) pprof() {
	s.AddRouter("/debug/pprof/cmdline", "GET", http.HandlerFunc(pprof.Cmdline))
	s.AddRouter("/debug/pprof/profile", "GET", http.HandlerFunc(pprof.Profile))
	s.AddRouter("/debug/pprof/heap", "GET", pprof.Handler("heap"))
	s.AddRouter("/debug/pprof/symbol", "GET", http.HandlerFunc(pprof.Symbol))
}

func (s *Server) Run(addr string) {

	// mux := http.NewServeMux()
	if s.profiler {
		s.pprof()
	}
	// mux.Handle("/", s)
	s.Logger.Printf("server serving %s\n", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.Logger.Fatal("ListenAndServe Fatal:", err)
	}
	s.l = l
	err = http.Serve(s.l, s)
	s.l.Close()
}
