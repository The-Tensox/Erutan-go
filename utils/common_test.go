package utils

import (
	"testing"
)

func TestSpaceComponent_RandomPositionInsideCircle(t *testing.T) {
	radius := 100.0
	for i := 0; i < 20; i++ {
		p := RandomPositionInsideCircle(radius)
		if p.X < -radius || p.X > radius || p.Z < -radius || p.Z > radius {
			t.Error("Incorrect RandomPositionInsideCircle")
		}
	}
}
