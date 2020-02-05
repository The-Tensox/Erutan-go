package game

import (
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"

	"github.com/user/erutan/ecs"
)

// AddLife set health component life, clip it and return true if entity is dead
func AddLife(h *erutan.Component_HealthComponent, value float64) bool {
	h.Life += value
	// Clip 0, 100
	if h.Life > 100 {
		h.Life = 100
	} else if h.Life < 0 {
		h.Life = 0
	}
	if h.Life == 0 {
		return true
	}
	return false
}

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
	*erutan.Component_HealthComponent
}

type ReachTargetSystem struct {
	entities []reachTargetEntity
}

func (r *ReachTargetSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	target *erutan.Component_TargetComponent,
	health *erutan.Component_HealthComponent) {
	r.entities = append(r.entities, reachTargetEntity{basic, space, target, health})
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
		if AddLife(entity.Component_HealthComponent, -10*dt) {
			utils.DebugLogf("I died, %v", entity.ID())
			r.Remove(*entity.BasicEntity)
		}
		// If I don't have a target, let's find one
		if entity.Target == nil {
			for _, e := range ManagerInstance.World.Systems() {
				if f, ok := e.(*EatableSystem); ok {
					//utils.DebugLogf("I'm %v, found a target: %v", entity.Component_SpaceComponent.Position, f.entities[0].Position)
					entity.Target = f.entities[0].Position // TODO: atm will just rush the first food of the array
					// Maybe later could finc the nearest, w/e ..
				}
			}
		}
		if entity.Target == nil {
			continue // There is no eatable ?
		}
		distance := utils.Distance(*entity.Position, *entity.Target)
		speed := 20.0
		newPos := utils.Add(*entity.Position,
			utils.Mul(utils.Div(utils.Sub(*entity.Target, *entity.Position), distance), dt*speed))
		//utils.DebugLogf("newpos %v", newPos)

		entity.Position = &newPos
	}
}

func (r *ReachTargetSystem) NotifyCallback(event utils.Event) {
	switch event.EventID {
	case utils.EntitiesCollided:
		for _, entity := range r.entities {
			entity.Target = nil // Find a new target
			AddLife(entity.Component_HealthComponent, 20)
			//utils.DebugLogf("my life %v", entity.Life)
		}
	}
}
