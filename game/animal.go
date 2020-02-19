package game

import (
	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"
	"github.com/user/erutan/utils/vector"

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
	*erutan.Component_PhysicsComponent
}

type herbivorousEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	Target *AnyObject
	*erutan.Component_HealthComponent
	*erutan.Component_SpeedComponent
}

type HerbivorousSystem struct {
	entities []herbivorousEntity
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
		volume := entity.Component_SpaceComponent.Scale.X * entity.Component_SpaceComponent.Scale.Y * entity.Component_SpaceComponent.Scale.Z

		// Every animal lose life proportional to deltatime, volume and speed
		// So bigger and faster animals need more food
		if AddLife(entity.Component_HealthComponent, -30*dt*volume*(entity.MoveSpeed/100)) {
			ManagerInstance.Broadcast <- erutan.Packet{
				Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
				Type: &erutan.Packet_DestroyEntity{
					DestroyEntity: &erutan.Packet_DestroyEntityPacket{
						EntityId: entity.ID(),
					},
				},
			}
			ManagerInstance.World.RemoveEntity(*entity.BasicEntity)
		}

		// If I don't have a target, let's find one
		if entity.Target == nil {
			h.findTarget(indexEntity, &entity)
		}
		if entity.Target == nil {
			continue // There is no target / animals
		}
		distance := vector.Distance(*entity.Position, *entity.Target.Position)
		newPos := vector.Add(*entity.Position,
			vector.Mul(vector.Div(vector.Sub(*entity.Target.Position, *entity.Position), distance), dt*entity.MoveSpeed))
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
		if f, ok := e.(*EatableSystem); ok && len(f.entities) > 0 {
			minPosition := f.entities[0]
			for _, eatableEntity := range f.entities {
				if vector.Distance(*entity.Position, *eatableEntity.Position) < vector.Distance(*entity.Position, *minPosition.Position) {
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
				speed := ((a.MoveSpeed + b.MoveSpeed) / 2) * utils.RandFloats(0.5, 1.5)
				scale := vector.Mul(vector.Div(vector.Add(*a.Scale, *b.Scale), 2), utils.RandFloats(0.5, 1.5))

				// Clipping scale ... TODO: add min & max scale somewhere
				clip := func(val float64, min float64, max float64) float64 {
					if val < 0.1 {
						return 0.1
					} else if val > 5 {
						return 5
					} else {
						return val
					}
				}
				scale.X = clip(scale.X, 0.1, 5)
				scale.Y = clip(scale.Y, 0.1, 5)
				scale.Z = clip(scale.Z, 0.1, 5)
				speed = clip(speed, 5, 80)
				position := a.Position
				position.Y = scale.Y // To stay above ground
				//utils.DebugLogf("Scale: %v", scale)
				ManagerInstance.AddHerbivorous(position, &scale, speed)
			}
		}
	}
}
