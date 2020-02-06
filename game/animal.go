package game

import (
	"math"

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
	erutan.Component_BehaviourTypeComponent
	erutan.Component_SpeedComponent
}

type reachTargetEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	*erutan.Component_TargetComponent
	*erutan.Component_HealthComponent
	*erutan.Component_SpeedComponent
}

type ReachTargetSystem struct {
	entities []reachTargetEntity
}

func (r *ReachTargetSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	target *erutan.Component_TargetComponent,
	health *erutan.Component_HealthComponent,
	speed *erutan.Component_SpeedComponent) {
	r.entities = append(r.entities, reachTargetEntity{basic, space, target, health, speed})
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
		if AddLife(entity.Component_HealthComponent, -3*dt) {
			utils.DebugLogf("I died, %v", entity.ID())
			r.Remove(*entity.BasicEntity)
		}
		// If I don't have a target, let's find one
		if entity.Target == nil {
			for _, e := range ManagerInstance.World.Systems() {
				if f, ok := e.(*EatableSystem); ok {
					//utils.DebugLogf("I'm %v, found a target: %v", entity.Component_SpaceComponent.Position, f.entities[0].Position)
					min := &erutan.NetVector3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
					for _, eatableEntity := range f.entities {
						if utils.Distance(*entity.Position, *eatableEntity.Position) < utils.Distance(*entity.Position, *min) {
							min = eatableEntity.Position
						}
					}
					entity.Target = min
				}
			}
		}
		if entity.Target == nil {
			continue // There is no eatable ?
		}
		distance := utils.Distance(*entity.Position, *entity.Target)
		newPos := utils.Add(*entity.Position,
			utils.Mul(utils.Div(utils.Sub(*entity.Target, *entity.Position), distance), dt*entity.MoveSpeed))
		//utils.DebugLogf("newpos %v", newPos)

		entity.Position = &newPos
	}
}

func (r *ReachTargetSystem) NotifyCallback(event utils.Event) {
	switch e := event.Value.(type) {
	case EntitiesCollided:
		for _, entity := range r.entities {
			entity.Target = nil // Find a new target
			AddLife(entity.Component_HealthComponent, 20*e.dt)
			//utils.DebugLogf("my life %v", entity.Life)
		}
	}
}

type animalReproductionEntity struct {
	*ecs.BasicEntity
	me     *Herbivorous
	target *Herbivorous
}

type AnimalReproductionSystem struct {
	entities []animalReproductionEntity
}

func (a *AnimalReproductionSystem) Add(basic *ecs.BasicEntity,
	me *Herbivorous,
	target *Herbivorous) {
	a.entities = append(a.entities, animalReproductionEntity{basic, me, target})
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (a *AnimalReproductionSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range a.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		a.entities = append(a.entities[:delete], a.entities[delete+1:]...)
	}
}

func (a *AnimalReproductionSystem) Update(dt float64) {
	var first animalReproductionEntity
	for _, entity := range a.entities {
		if entity.me.Life > 80 {
			if first == (animalReproductionEntity{}) {
				first = entity
			} else {
				//utils.DebugLogf("Reproduction mode")
				entity.target = first.me
				first.target = entity.me
				return
			}
		}
		if entity.target != (&Herbivorous{}) {
			entity.me.Target = entity.target.Position
		}
	}
}

func (a *AnimalReproductionSystem) NotifyCallback(event utils.Event) {
	switch u := event.Value.(type) {
	case EntitiesCollided:
		// If an animal collided with me
		if u.a.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			u.b.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL {
			// utils.DebugLogf("Reproduction")
		}
	}
}
