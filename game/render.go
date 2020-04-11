package game

import (
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/octree"
)

type renderObject struct {
	Id uint64
	*erutan.Component_RenderComponent
}

type RenderSystem struct {
	objects octree.Octree
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
	/*
		for _, entity := range r.entities {

		}
	*/
}
