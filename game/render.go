package game

import (
	"github.com/The-Tensox/erutan/cfg"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

type renderObject struct {
	Id uint64
	*erutan.Component_RenderComponent
}

type RenderSystem struct {
	objects octree.Octree
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(),
		cfg.Global.Logic.GroundSize*1000))}
}

func (r *RenderSystem) Add(id uint64,
	render *erutan.Component_RenderComponent) {
	// Non-physical objects have default position 0,0,0
	r.objects.Insert(*octree.NewObjectCube(renderObject{id, render}, 0, 0, 0, 1))
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (r *RenderSystem) Remove(object octree.Object) {
	r.objects.Remove(object)
}

func (r *RenderSystem) Update(dt float64) {
	//o := r.objects.GetAllObjects()
	//for i := range o {
	//	if ro, ok := o[i].Data.(renderObject); ok && ro.Green > 10{
	//		utils.DebugLogf("ground %v", ro.Id)
	//	}
	//}
}
