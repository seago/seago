package cache

import (
	"sync"
)

type Storer interface {
	Get(key interface{}) interface{}
	Set(key, value interface{})
	Len() int
	Delete(key interface{}) interface{}
}

type Cache struct {
	mux   *sync.RWMutex
	store Storer
}

func NewCache(store Storer) *Cache {
	return &Cache{
		mux:   new(sync.RWMutex),
		store: store,
	}
}

func (c *Cache) Get(key interface{}) interface{} {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.store.Get(key)
}

func (c *Cache) Set(key, value interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.store.Set(key, value)
}

func (c *Cache) Delete(key interface{}) interface{} {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.store.Delete(key)
}
