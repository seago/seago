package middleware

import (
	"testing"
)

func TestGet(t *testing.T) {
	DefaultMiddleware.Add("db", "mysql")
	db := DefaultMiddleware.Get("db").(string)
	if db != "mysql" {
		t.Error("middleware get error")
	}
}

func TestSet(t *testing.T) {
	DefaultMiddleware.Add("db", "mysql")
	DefaultMiddleware.Set("db", "sql server")
	db := DefaultMiddleware.Get("db").(string)
	if db == "mysql" {
		t.Error("middleware set error")
	}
}
