package vector

import (
	"testing"

	erutan "github.com/user/erutan/protos/realtime"
)

func TestNewBox(t *testing.T) {
	b := NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(-1, -1, -1))
	if &b.Min == erutan.NewNetVector3(-1, -1, -1) && &b.Max == erutan.NewNetVector3(0, 0, 0) {
		t.Fail()
	}

	b = NewBox(erutan.NewNetVector3(0, 0, 0), erutan.NewNetVector3(1, 1, 1))
	if &b.Min == erutan.NewNetVector3(0, 0, 0) && &b.Max == erutan.NewNetVector3(1, 1, 1) {
		t.Fail()
	}
}
