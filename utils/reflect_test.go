package utils

import (
	"changit/http_rpc/player"
	"testing"
)

func TestGetMethods(t *testing.T) {
	GetMethods(player.DefaultPlayerManager)
}

type ATest struct {
	B string
	C int
	BTest
}

type BTest struct {
	D string
	E float64
}

func TestIsEmptyValue(t *testing.T) {
	a1 := ATest{}
	a2 := ATest{B: "a2"}
	t.Log(IsEmptyValue(a1))
	t.Log(IsEmptyValue(a2))

}
