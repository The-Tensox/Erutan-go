package utils

import (
	"testing"

	erutan "github.com/user/erutan/protos/realtime"
)

func TestSpaceComponent_RandomPositionInsideCircle(t *testing.T) {
	radius := 100.0
	for i := 0; i < 20; i++ {
		p := RandomPositionInsideCircle(&erutan.NetVector2{X: 0, Y: 0}, radius)
		if p.X < -radius || p.X > radius || p.Z < -radius || p.Z > radius {
			t.Error("Incorrect RandomPositionInsideCircle")
		}
	}
}
