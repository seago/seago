package context

import (
	"testing"
)

func TestString(t *testing.T) {
	value := Value("123")
	v := value.String()
	if v != "123" {
		t.Fail()
	}
}

func TestInt(t *testing.T) {
	value := Value("123")
	v := value.Int()
	if v != 123 {
		t.Fail()
	}
}

func TestInt32(t *testing.T) {
	value := Value("123")
	v := value.Int32()
	if v != int32(123) {
		t.Fail()
	}
}

func TestInt64(t *testing.T) {
	value := Value("123")
	v := value.Int64()
	if v != int64(123) {
		t.Fail()
	}
}

func TestFloat32(t *testing.T) {
	value := Value("123")
	v := value.Float32()
	if v != float32(123) {
		t.Fail()
	}
}

func TestFloat64(t *testing.T) {
	value := Value("123")
	v := value.Float64()
	if v != float64(123) {
		t.Fail()
	}
}

func TestJsonDecode(t *testing.T) {
	value := Value(`{"key":"key1","value":"value1"}`)
	type test struct {
		Key   string `json:key`
		Value string `json:value`
	}
	tt := &test{}
	err := value.JsonDecode(tt)
	if err != nil {
		t.Error(err)
	}

	if tt.Key != "key1" || tt.Value != "value1" {
		t.Error("json decode error")
	}
}
