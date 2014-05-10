package router

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAddRouter(t *testing.T) {
	router := NewRouter()
	handler := func() string {
		return "test router"
	}
	err := router.AddRouter("/test/:name", "GET", handler)
	if err != nil {
		t.Error(err)
	}

	handler1 := "123"
	err = router.AddRouter("/test", "GET", handler1)
	if err == nil {
		t.Fail()
	}

	handler2 := func(wr http.ResponseWriter, r *http.Request) {
		wr.Write([]byte("test http hanlerfunc"))
	}
	err = router.AddRouter("/test", "GET", handler2)
	if err != nil {
		t.Error(err)
	}
}
