package memory

import (
	"container/list"
)

type MemoryCache struct {
	maxEntries int
	ll         *list.List
	cache      map[interface{}]*list.Element
}

func New(maxEntries int) *MemoryCache {
	return &MemoryCache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element, maxEntries),
	}
}

type entry struct {
	key   interface{}
	value interface{}
}

func (m *MemoryCache) Get(key interface{}) interface{} {
	if m.cache == nil {
		return nil
	}
	if elem, ok := m.cache[key]; ok {
		m.ll.MoveToFront(elem)
		return elem.Value.(*entry).value
	}
	return nil
}

func (m *MemoryCache) Set(key, value interface{}) {
	if m.cache == nil {
		m.cache = make(map[interface{}]*list.Element)
		m.ll = list.New()
	}
	if elem, ok := m.cache[key]; ok {
		m.ll.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}
	elem := &entry{key, value}
	e := m.ll.PushFront(elem)
	m.cache[key] = e
	if m.maxEntries != 0 && m.ll.Len() > m.maxEntries {
		m.deleteOldEntries()
	}
}

func (m *MemoryCache) Delete(key interface{}) interface{} {
	if m.cache == nil {
		return nil
	}
	if elem, ok := m.cache[key]; ok {
		m.deleteElem(elem)
		return elem.Value.(*entry).value
	}
	return nil
}

func (m *MemoryCache) deleteOldEntries() {
	if m.cache == nil {
		return
	}
	elem := m.ll.Back()
	if elem != nil {
		m.deleteElem(elem)
	}
}

func (m *MemoryCache) deleteElem(e *list.Element) {
	m.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(m.cache, kv.key)
}

func (m *MemoryCache) Len() int {
	if m.cache == nil {
		return 0
	}
	return m.ll.Len()
}
