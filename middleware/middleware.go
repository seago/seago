package middleware

import (
	"sync"
)

var DefaultMiddleware = NewMiddleware()

type Middleware struct {
	lock *sync.Mutex
	data map[string]interface{}
}

func NewMiddleware() *Middleware {
	return &Middleware{lock: new(sync.Mutex), data: make(map[string]interface{})}
}

func (m *Middleware) Get(key string) interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.data[key]
}

func (m *Middleware) Set(key string, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[key] = value
}

func (m *Middleware) Add(key string, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[key] = value
}
