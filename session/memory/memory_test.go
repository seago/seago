package memory

import (
	"testing"
	"time"
)

type session struct {
	name string
}

var pass = session{"123"}
var miss = session{"456"}

var testSession = map[string]interface{}{
	"test_pass":   123456,
	"test_miss":   12345,
	"struct_pass": pass,
	"struct_miss": miss,
}

func TestGet(t *testing.T) {
	ms := New(20)
	for k, v := range testSession {
		ms.Set(k, v, 0)
	}
	if ms.Get("test_pass").(int) != 123456 {
		t.Error("memory session get error")
	}

	if ms.Get("test_miss").(int) == 123456 {
		t.Error("memory session get errorr")
	}

	if ms.Get("struct_pass").(session) != pass {
		t.Error("memory session get error")
	}
	if ms.Get("struct_miss").(session) == pass {
		t.Error("memory session get error")
	}
}

func TestExpires(t *testing.T) {
	ms := New(20)
	for k, v := range testSession {
		ms.Set(k, v, 0)
	}
	time.Sleep(21 * time.Second)
	if !ms.Expires() {
		t.Error("Session not expires")
	}
}

func TestStorageExpires(t *testing.T) {
	ms := New(20)
	ms.Set("test", "123", 10)
	time.Sleep(11 * time.Second)
	v := ms.Get("test")
	if v != nil {
		t.Error("storage expires error")
	}
	if ms.storge["test"].expires > time.Now().Unix() {
		t.Error("storage is not expires")
	}
	ms.Set("test1", "123", 20)
	time.Sleep(11 * time.Second)
	v = ms.Get("key1")
	if v != nil {
		t.Error("storage expires error")
	}
	if ms.storge["test1"].expires > time.Now().Unix() {
		t.Error("storage is not expires")
	}
}

func TestGC(t *testing.T) {
	ms := New(20)
	ms.Set("test", "123", 10)
	time.Sleep(11 * time.Second)
	ms.GC()
	v := ms.Get("test")
	if v != nil {
		t.Error("GC error")
	}
	time.Sleep(11 * time.Second)
	ms.GC()
	if len(ms.storge) != 0 {
		t.Error("GC error")
	}
}
