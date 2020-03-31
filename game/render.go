package game

import (
	erutan "github.com/The-Tensox/erutan/protobuf"

	"github.com/The-Tensox/erutan/ecs"
)

type renderEntity struct {
	*ecs.BasicEntity
	*erutan.Component_RenderComponent
}

type RenderSystem struct {
	entities []renderEntity
}

func (r *RenderSystem) Add(basic *ecs.BasicEntity,
	render *erutan.Component_RenderComponent) {
	r.entities = append(r.entities, renderEntity{basic, render})
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (r *RenderSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range r.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		r.entities = append(r.entities[:delete], r.entities[delete+1:]...)
	}
}

func (r *RenderSystem) Update(dt float64) {
	/*
		for _, entity := range r.entities {

		}
	*/
}
