package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

type Storager interface {
	Get(key string) interface{}
	Set(key string, value interface{}, exprise int64) error
	SetExpires(expires int64)
	Clear(key string)
	Flush()
	Expires() bool
	GC()
}

type SessionManager struct {
	storage     map[string]Storager
	Name        string //name of cookie
	Domain      string //domain of cookie
	Path        string //path of cookie
	HttpOnly    bool   //httpOnly of cookie
	Secure      bool   // secure of cookie
	GcDelay     time.Duration
	storageName string
	expires     int64
}

func New(storageName string, expires int64) *SessionManager {
	sm := &SessionManager{
		storage:     make(map[string]Storager),
		Name:        "SeagoSID",
		Path:        "/",
		Domain:      "",
		HttpOnly:    true,
		Secure:      true,
		GcDelay:     time.Minute * 30,
		storageName: storageName,
		expires:     expires,
	}
	return sm
}

func (sm *SessionManager) NewSession(w http.ResponseWriter) Storager {
	sid, err := generate_sid()
	if err != nil {
		panic("Session:NewSession error")
	}
	cookie := &http.Cookie{}
	cookie.Value = sid
	cookie.Name = sm.Name
	cookie.Path = sm.Path
	cookie.Domain = sm.Domain
	cookie.HttpOnly = sm.HttpOnly
	cookie.Secure = sm.Secure
	storages[sm.storageName].SetExpires(sm.expires) //set storager expires

	sm.storage[sid] = storages[sm.storageName]

	http.SetCookie(w, cookie)
	return sm.storage[sid]
}

func (sm *SessionManager) GetStorager(w http.ResponseWriter, r *http.Request) Storager {
	cookie, err := r.Cookie(sm.Name)
	if err != nil {
		return sm.NewSession(w)
	}

	if _, ok := sm.storage[cookie.Value]; !ok {
		return sm.NewSession(w)
	}

	return sm.storage[cookie.Value]
}

func (sm *SessionManager) GC() {
	for {
		select {
		case <-time.After(sm.GcDelay):
			for k, v := range sm.storage {
				if v.Expires() {
					delete(sm.storage, k)
				} else {
					v.GC()
				}
			}
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
