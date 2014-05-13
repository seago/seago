package memory

import (
	"errors"
	"github.com/seago/seago/session"
	"sync"
	"time"
)

type storage struct {
	value   interface{}
	expires int64
}

func newStorage(value interface{}, expires int64) *storage {
	return &storage{value, expires}
}

type MemorySession struct {
	storge  map[string]*storage
	expires int64
	*sync.Mutex
}

var (
	setError = errors.New("session expires later than storage expires")
)

func New(expires int64) *MemorySession {
	return &MemorySession{make(map[string]*storage), expires + time.Now().Unix(), new(sync.Mutex)}
}

func (ms *MemorySession) Get(key string) interface{} {
	ms.Lock()
	defer ms.Unlock()
	if ms.Expires() || ms.storageExpires(key) {
		return nil
	}
	if _, ok := ms.storge[key]; !ok {
		return nil
	}
	return ms.storge[key].value
}

func (ms *MemorySession) Set(key string, value interface{}, expires int64) error {
	ms.Lock()
	defer ms.Unlock()
	if ms.Expires() {
		return setError
	}
	if expires == 0 || time.Now().Unix()+expires > ms.expires {
		expires = ms.expires
	} else {
		expires = time.Now().Unix() + expires
	}
	st := newStorage(value, expires)
	ms.storge[key] = st
	return nil
}

func (ms *MemorySession) SetExpires(expired int64) {
	ms.expires = expired + time.Now().Unix()
}

func (ms *MemorySession) Clear(key string) {
	ms.Lock()
	defer ms.Unlock()
	if _, ok := ms.storge[key]; ok {
		delete(ms.storge, key)
	}
}

func (ms *MemorySession) Flush() {
	ms = nil
}

func (ms *MemorySession) storageExpires(key string) bool {
	now := time.Now().Unix()
	if ms.Expires() {
		return true
	}
	if _, ok := ms.storge[key]; !ok {
		return true
	}
	if ms.storge[key].expires < now {
		return true
	}
	return false
}

func (ms *MemorySession) Expires() bool {
	now := time.Now().Unix()
	if ms.expires < now {
		return true
	}
	return false
}

func (ms *MemorySession) GC() {
	if ms.Expires() {
		ms = nil
		return
	}
	for k, _ := range ms.storge {
		if ms.storageExpires(k) {
			delete(ms.storge, k)
		}
	}
}

func init() {
	session.Register("memory", New(604800))
}
