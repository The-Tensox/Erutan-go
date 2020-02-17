package game

import (
	"testing"

	erutan "github.com/user/erutan/protos/realtime"
)

func TestSpaceComponent_Collision(t *testing.T) {
	componentA := &erutan.Component_SpaceComponent{
		Position: &erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
	}
	aabb := GetAABB(componentA)
	t.Log("AABB", aabb)
	solution := AABB{minX: -0.5, maxX: 0.5, minY: -0.5, maxY: 0.5, minZ: -0.5, maxZ: 0.5}
	if aabb != solution {
		t.Errorf("Should be equal to %v", solution)
	}

	componentB := &erutan.Component_SpaceComponent{
		Position: &erutan.NetVector3{X: 0, Y: 0, Z: 0},
		Rotation: &erutan.NetQuaternion{X: 0, Y: 0, Z: 0, W: 0},
		Scale:    &erutan.NetVector3{X: 1, Y: 1, Z: 1},
	}
	t.Log("Overlap", Overlap(componentA, componentB))

}
