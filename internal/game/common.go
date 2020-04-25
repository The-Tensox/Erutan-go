package game

import (
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

// Move to octree package ?
func Find(oc octree.Octree, object octree.Object) *octree.Object {
	size := object.Bounds.GetSize() // atm assuming cube
	// HACK to increase a little bit the box otherwise doesn't detect
	newMin := object.Bounds.Min.Minus(size.Times(0.1))
	newMax := object.Bounds.Max.Plus(size.Times(0.1))
	newBox := protometry.Box{
		Min: &newMin,
		Max: &newMax,
	}
	objs := oc.GetColliding(newBox)
	// We need to find the current Object in collisionSystem's Octree
	for i := range objs {
		if objs[i].Equal(object) { // Could instead compare ids
			return &objs[i]
		}
	}
	//utils.DebugLogf("nil")
	//for _, o := range oc.GetAllObjects() {
	//	utils.DebugLogf("%v", o.ID())
	//}
	return nil
}