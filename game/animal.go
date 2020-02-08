package game

import (
	"math/rand"

	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"

	"github.com/user/erutan/ecs"
)

// AddLife set health component life, clip it and return true if entity is dead
func AddLife(h *erutan.Component_HealthComponent, value float64) bool {
	h.Life += value
	// Clip 0, 100
	/*if h.Life > 100 {
		h.Life = 100
	} else*/if h.Life < 0 {
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

type herbivorousEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	Target *AnyObject
	*erutan.Component_HealthComponent
	*erutan.Component_SpeedComponent
}

type HerbivorousSystem struct {
	entities        []herbivorousEntity
	speedStatistics erutan.Statistics // Keep track of stats
	lifeStatistics  erutan.Statistics
}

func (h *HerbivorousSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent,
	target *AnyObject,
	health *erutan.Component_HealthComponent,
	speed *erutan.Component_SpeedComponent) {
	h.entities = append(h.entities, herbivorousEntity{basic, space, target, health, speed})
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (h *HerbivorousSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range h.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		h.entities = append(h.entities[:delete], h.entities[delete+1:]...)
	}
}

func (h *HerbivorousSystem) Update(dt float64) {
	for indexEntity, entity := range h.entities {
		h.computeStatistics(entity.MoveSpeed, entity.Life)
		if AddLife(entity.Component_HealthComponent, -3*dt) {
			ManagerInstance.World.RemoveEntity(*entity.BasicEntity)
		}

		// If I don't have a target, let's find one
		if entity.Target == nil {
			h.findTarget(indexEntity, &entity)
		}
		if entity.Target == nil {
			continue // There is no target / animals
		}
		distance := utils.Distance(*entity.Position, *entity.Target.Position)
		newPos := utils.Add(*entity.Position,
			utils.Mul(utils.Div(utils.Sub(*entity.Target.Position, *entity.Position), distance), dt*entity.MoveSpeed))
		entity.Position = &newPos
	}
}

func (h *HerbivorousSystem) findTarget(indexEntity int, entity *herbivorousEntity) {
	if entity.Life > 80 {
		for j := indexEntity + 1; j < len(h.entities); j++ {
			if h.entities[j].Life > 80 {
				newTarget := &AnyObject{}
				newTarget.BasicEntity = h.entities[j].BasicEntity
				newTarget.Component_SpaceComponent = h.entities[j].Component_SpaceComponent
				entity.Target = newTarget

				newTargetTwo := &AnyObject{}
				newTargetTwo.BasicEntity = entity.BasicEntity
				newTargetTwo.Component_SpaceComponent = entity.Component_SpaceComponent
				h.entities[j].Target = newTargetTwo
				// Found a target (another animal)
				return
			}
		}
	}
	for _, e := range ManagerInstance.World.Systems() {
		if f, ok := e.(*EatableSystem); ok {
			minPosition := f.entities[0]
			for _, eatableEntity := range f.entities {
				if utils.Distance(*entity.Position, *eatableEntity.Position) < utils.Distance(*entity.Position, *minPosition.Position) {
					minPosition = eatableEntity
				}
			}
			newTarget := &AnyObject{}
			newTarget.Component_SpaceComponent = minPosition.Component_SpaceComponent
			newTarget.BasicEntity = minPosition.BasicEntity
			entity.Target = newTarget
		}
	}
}

func (h *HerbivorousSystem) Find(id uint64) *herbivorousEntity {
	for _, entity := range h.entities {
		if entity.ID() == id {
			return &entity
		}
	}
	return nil
}
func (h *HerbivorousSystem) NotifyCallback(event utils.Event) {
	switch e := event.Value.(type) {
	case EntitiesCollided:
		a := h.Find(e.a.ID())
		b := h.Find(e.b.ID())
		if e.a.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION &&
			e.b.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL {
			AddLife(b.Component_HealthComponent, 40)

			// Reset target for everyone that had this target
			for _, e := range h.entities {
				if (e.Target != nil && a != nil) && (e.Target.ID() == a.ID()) {
					e.Target = nil
				}
			}
		} else if e.b.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION &&
			e.a.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL {
			AddLife(b.Component_HealthComponent, 40)

			// Reset target for everyone that had this target
			for _, e := range h.entities {
				if (e.Target != nil && b != nil) && (e.Target.ID() == b.ID()) {
					e.Target = nil
				}
			}
		} else { // Both are animals
			if (a != nil && b != nil) && (a.Life > 80 && b.Life > 80) {
				if a.Target != nil {
					a.Target = nil
				}
				if b.Target != nil {
					b.Target = nil
				}
				AddLife(a.Component_HealthComponent, -50)
				AddLife(b.Component_HealthComponent, -50)
				speed := ((a.MoveSpeed + b.MoveSpeed) / 2) * (1 + (-0.5 + rand.Float64()*1))
				ManagerInstance.AddHerbivorous(a.Position, speed)
			}
		}
	}
}

func (h *HerbivorousSystem) computeStatistics(speed float64, life float64) {
	if speed != -1 {
		h.speedStatistics.Average = (h.speedStatistics.Average + speed) / 2
		if speed < h.speedStatistics.Minimum {
			h.speedStatistics.Minimum = speed
		}
		if speed > h.speedStatistics.Maximum {
			h.speedStatistics.Maximum = speed
		}
	}

	if life != -1 {
		h.lifeStatistics.Average = (h.lifeStatistics.Average + life) / 2
		if life < h.lifeStatistics.Minimum {
			h.lifeStatistics.Minimum = life
		}
		if life > h.lifeStatistics.Maximum {
			h.lifeStatistics.Maximum = life
		}
	}
	// utils.DebugLogf("Statistics %v - %v", h.speedStatistics, h.lifeStatistics)
}
