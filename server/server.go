package server

import (
	"github.com/seago/seago/router"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
)

var VERSION = "0.0.1"

type Server struct {
	router    *router.Router
	Logger    *log.Logger
	l         net.Listener
	profiler  bool
	maxMemory int64
	version   string
}

func NewServer() *Server {
	return &Server{
		router:    router.NewRouter(),
		Logger:    log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
		maxMemory: 100 << 20, //100M
	}
}

func (s *Server) SetProfile(is bool) {
	s.profiler = is
}

func (s *Server) SetMaxMemory(maxMemory int64) {
	s.maxMemory = maxMemory
}

func (s *Server) AddRouter(pattern, method string, handler interface{}) {
	err := s.router.AddRouter(pattern, method, handler)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Server", "seago "+VERSION)
	s.router.Process(rw, r, s.maxMemory)
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
