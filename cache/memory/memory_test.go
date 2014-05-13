package memory

import (
	"strings"
	"testing"
)

type test_struct1 struct {
	name string
	age  int
}

type test_struct2 struct {
	test_struct1
	num int
}

var test = []struct {
	name    string
	set_key interface{}
	get_key interface{}
}{
	{"string_hit", "myKey", "myKey"},
	{"string_miss", "myKey", "nokey"},
	{"struct1_hit", test_struct1{"miller", 12}, test_struct1{"miller", 12}},
	{"struct1_miss", test_struct1{"miller", 12}, test_struct1{"miller", 0}},
	{"struct2_hit", test_struct2{test_struct1{"john", 34}, 3}, test_struct2{test_struct1{"john", 34}, 3}},
	{"struct2_miss", test_struct2{test_struct1{"john", 34}, 3}, test_struct2{test_struct1{"jo1hn", 34}, 3}},
}

func TestGet(t *testing.T) {
	for _, v := range test {
		mc := New(0)
		mc.Set(v.set_key, 1234)
		val := mc.Get(v.get_key)
		if !strings.Contains(v.name, "miss") {
			if val == nil {
				t.Fatalf("%s: cache hit = %v; want %v", v.name, val, 1234)
			} else if val.(int) != 1234 {
				t.Fatalf("%s expected get to return 1234 but got %v", v.name, val)
			}
		} else {
			if val != nil {
				t.Fatalf("%s: cache hit = %v; want %v", v.name, val, nil)
			}
		}
	}
}

func TestDelete(t *testing.T) {
	mc := New(0)
	mc.Set("myKey", 1234)
	val := mc.Get("myKey")
	if val == nil {
		t.Fatal("TestDelete returned no match")
	} else if val != 1234 {
		t.Fatalf("TestDelete failed.  Expected %d, got %v", 1234, val)
	}

	mc.Delete("myKey")
	v := mc.Get("myKey")
	if v != nil {
		t.Fatal("TestDelete returned a delete entry")
	}
}
