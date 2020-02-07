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
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	*erutan.Component_HealthComponent
	Target *AnyObject
	*erutan.Component_RenderComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_SpeedComponent
}

type reachTargetEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	Target *AnyObject
	*erutan.Component_HealthComponent
	*erutan.Component_SpeedComponent
}

type ReachTargetSystem struct {
	entities []reachTargetEntity
}

func (r *ReachTargetSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	target *AnyObject,
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
	for indexEntity, entity := range r.entities {
		//utils.DebugLogf("my target %v", entity.Target)
		if AddLife(entity.Component_HealthComponent, -3*dt) {
			utils.DebugLogf("I died, %v", entity.ID())
			ManagerInstance.World.RemoveEntity(*entity.BasicEntity)
		}

		if entity.Life > 80 {
			for j := indexEntity + 1; j < len(r.entities); j++ {
				if r.entities[j].Life > 80 {
					newTarget := &AnyObject{}
					newTarget.BasicEntity = r.entities[j].BasicEntity
					newTarget.Component_SpaceComponent = r.entities[j].Component_SpaceComponent
					entity.Target = newTarget

					newTargetTwo := &AnyObject{}
					newTargetTwo.BasicEntity = entity.BasicEntity
					newTargetTwo.Component_SpaceComponent = entity.Component_SpaceComponent
					r.entities[j].Target = newTargetTwo
				}
			}
		}

		// If I don't have a target, let's find one
		if entity.Target == nil {
			//utils.DebugLogf("entitytarget %v", entity.Target)
			for _, e := range ManagerInstance.World.Systems() {
				if f, ok := e.(*EatableSystem); ok {
					minPosition := f.entities[0] //erutan.NetVector3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
					for _, eatableEntity := range f.entities {
						if utils.Distance(*entity.Position, *eatableEntity.Position) < utils.Distance(*entity.Position, *minPosition.Position) {
							minPosition = eatableEntity
						}
					}
					newTarget := &AnyObject{}
					newTarget.Component_SpaceComponent = minPosition.Component_SpaceComponent
					newTarget.BasicEntity = minPosition.BasicEntity
					entity.Target = newTarget
					//utils.DebugLogf("I'm %v, found a target: %v", entity.ID(), entity.Target)
				}
			}
		}
		if entity.Target == nil {
			continue // There is no eatable ?
		}
		//utils.DebugLogf("me %v %v %v", entity.ID(), entity.Position, entity.Target)
		distance := utils.Distance(*entity.Position, *entity.Target.Position)
		newPos := utils.Add(*entity.Position,
			utils.Mul(utils.Div(utils.Sub(*entity.Target.Position, *entity.Position), distance), dt*entity.MoveSpeed))
		//utils.DebugLogf("newpos %v", newPos)

		entity.Position = &newPos
	}
}

func (r *ReachTargetSystem) Find(id uint64) *reachTargetEntity {
	for _, entity := range r.entities {
		if entity.ID() == id {
			return &entity
		}
	}
	return nil
}
func (r *ReachTargetSystem) NotifyCallback(event utils.Event) {
	switch e := event.Value.(type) {
	case EntitiesCollided:
		a := r.Find(e.a.ID())
		b := r.Find(e.b.ID())
		if e.a.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION &&
			e.b.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL {
			AddLife(b.Component_HealthComponent, 20)

			// Reset target for everyone that had this target
			for _, e := range r.entities {
				if (e.Target != nil && a != nil) && (e.Target.ID() == a.ID()) {
					utils.DebugLogf("set nil")
					e.Target = nil
				}
			}
		} else if e.b.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION &&
			e.a.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL {
			AddLife(a.Component_HealthComponent, 20)

			// Reset target for everyone that had this target
			for _, e := range r.entities {
				if (e.Target != nil && b != nil) && (e.Target.ID() == b.ID()) {
					utils.DebugLogf("set nil")
					e.Target = nil
				}
			}
		} else {
			if (a != nil && b != nil) && (a.Life > 80 && b.Life > 80) {
				if a.Target != nil {
					a.Target = nil
				}
				if b.Target != nil {
					b.Target = nil
				}
				//utils.DebugLogf("Repro %v %v", a, b)
				AddLife(a.Component_HealthComponent, -50)
				AddLife(b.Component_HealthComponent, -50)
				//utils.DebugLogf("Repro %v %v", a.Life, b.Life)

			}
		}
		//}
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
	for i := 0; i < len(a.entities); i++ {
		if a.entities[i].me.Life > 80 {
			for j := i + 1; j < len(a.entities); j++ {
				if a.entities[j].me.Life > 80 {
					//utils.DebugLogf("i %v, %v, j %v, %v", i, a.entities[i].me.Life, j, a.entities[j].me.Life)
					//a.entities[i].target = a.entities[j].me
					//a.entities[j].target = a.entities[i].me
					newTarget := &AnyObject{}
					newTarget.BasicEntity = a.entities[j].me.BasicEntity
					newTarget.Component_SpaceComponent = a.entities[j].me.Component_SpaceComponent
					a.entities[i].me.Target = newTarget

					newTargetTwo := &AnyObject{}
					newTargetTwo.BasicEntity = a.entities[i].me.BasicEntity
					newTargetTwo.Component_SpaceComponent = a.entities[i].me.Component_SpaceComponent
					a.entities[j].me.Target = newTargetTwo
				}
			}
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
