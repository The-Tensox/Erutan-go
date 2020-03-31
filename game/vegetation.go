package game

import (
	"github.com/The-Tensox/erutan/ecs"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
)

type AnyObject struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	*erutan.Component_RenderComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_PhysicsComponent
}

type eatableEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
}

type EatableSystem struct {
	entities []eatableEntity
}

func (e *EatableSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent) {
	e.entities = append(e.entities, eatableEntity{basic, space})
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (e *EatableSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range e.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		e.entities = append(e.entities[:delete], e.entities[delete+1:]...)
	}
}

func (e *EatableSystem) Update(dt float64) {
	/*
		for _, entity := range e.entities {
		}
	*/
}

func (e *EatableSystem) NotifyCallback(event utils.Event) {
	switch u := event.Value.(type) {
	case EntitiesCollided:
		// If an animal collided with me
		if u.a.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			u.b.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			newSc := u.b.Component_SpaceComponent
			newSc.Position = utils.RandomPositionInsideCircle(&erutan.NetVector2{X: utils.Config.GroundSize / 2,
				Y: utils.Config.GroundSize / 2},
				utils.Config.GroundSize/2)
			ManagerInstance.Watch.Notify(utils.Event{Value: EntityPhysicsUpdated{id: u.b.ID(), newSc: *newSc, dt: u.dt}})

		}

		if u.b.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			u.a.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			newSc := u.a.Component_SpaceComponent
			newSc.Position = utils.RandomPositionInsideCircle(&erutan.NetVector2{X: utils.Config.GroundSize / 2,
				Y: utils.Config.GroundSize / 2},
				utils.Config.GroundSize/2)
			ManagerInstance.Watch.Notify(utils.Event{Value: EntityPhysicsUpdated{id: u.a.ID(), newSc: *newSc, dt: u.dt}})
		}

	}
}
