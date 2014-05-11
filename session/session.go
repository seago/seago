package session

import (
	"container/list"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

type Storager interface {
	Get(key string) interface{}
	Set(key string, value interface{}, exprise int64) error
	Clear(key string)
	Flush()
	Expires() bool
	GC()
}

type SessionManager struct {
	lock        *sync.Mutex
	storage     map[string]Storager
	cookieName  string
	path        string
	httpOnly    bool
	secure      bool
	gcDelay     int64
	storageName string
}

var storages = make(map[string]Storager)

func Register(name string, storage Storager) {
	if storage == nil {
		panic("Session: Register storage is nil")
	}
	if _, ok := storages[name]; ok {
		panic("Session: Storge is registed")
	}
	storages[name] = storage
}

func New() *SessionManager {
	sm := &SessionManager{
		lock:        new(sync.Mutex),
		stroage:     make(map[string]Storager),
		cookieName:  "SeagoSID",
		path:        "/",
		httpOnly:    true,
		secure:      true,
		gcDelay:     time.Minute * 30,
		storageName: "memory",
	}
}

func (sm *SessionManager) NewSession(w http.ResponseWriter) Storager {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sid := generate_sid()
	cookie := &http.Cookie{}
	cookie.Value = sid
	cookie.Name = sm.cookieName
	cookie.Path = sm.path
	cookie.HttpOnly = sm.httpOnly
	cookie.Secure = sm.secure
	sm.storage[sid] = storages[sm.storageName]
	http.SetCookie(w, cookie)
	return sm.storage[sid]
}

func (sm *SessionManager) GetStorager(w http.ResponseWriter, r *http.Request) Storager {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return sm.NewSession(w)
		}
		return nil
	}
	return sm.storage[cookie.Value]
}

func (sm *SessionManager) GC() {
	for {
		select {
		case <-time.After(sm.gcDelay):
			sm.lock.Lock()
			for k, v := range sm.storage {
				if v.Expires() {
					delete(sm.storage, k)
				} else {
					v.GC()
				}
			}
			sm.lock.Unlock()
		}
	}
}

func generate_sid() (sid string, err error) {
	// Following code from: http://www.ashishbanerjee.com/home/go/go-generate-uuid
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7
	sid = hex.EncodeToString(uuid)
	return
}
