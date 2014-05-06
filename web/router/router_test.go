package router

import (
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
}
