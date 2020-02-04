package game

import (
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"

	"github.com/user/erutan/ecs"
)

type Herbivorous struct {
	ecs.BasicEntity
	erutan.Component_SpaceComponent
	erutan.Component_HealthComponent
	erutan.Component_TargetComponent
	erutan.Component_RenderComponent
}

type reachTargetEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	*erutan.Component_TargetComponent
}

type ReachTargetSystem struct {
	entities []reachTargetEntity
}

func (r *ReachTargetSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	target *erutan.Component_TargetComponent) {
	r.entities = append(r.entities, reachTargetEntity{basic, space, target})
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (r *ReachTargetSystem) Remove(basic ecs.BasicEntity) {
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

func (r *ReachTargetSystem) Update(dt float64) {
	for _, entity := range r.entities {
		// If I don't have a target, let's find one
		if entity.Target == nil {
			for _, e := range ManagerInstance.World.Systems() {
				if f, ok := e.(*EatableSystem); ok {
					utils.DebugLogf("I'm %v, found a target: %v", entity.Component_SpaceComponent.Position, f.entities[0].Position)
					entity.Target = f.entities[0].Position // TODO: atm will just rush the first food of the array
					// Maybe later could finc the nearest, w/e ..
				}
			}
		}
		if entity.Target == nil {
			continue // TODO: or move random
		}
		distance := utils.Distance(*entity.Position, *entity.Target)
		newPos := utils.Add(*entity.Position,
			utils.Div(utils.Sub(*entity.Target, *entity.Position), dt*distance*10))
		entity.Position = &newPos
	}
}
