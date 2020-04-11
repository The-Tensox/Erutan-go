package game

import (
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

type AnyObject struct {
	Id uint64
	*erutan.Component_SpaceComponent
	*erutan.Component_RenderComponent
	*erutan.Component_BehaviourTypeComponent
	*erutan.Component_PhysicsComponent
}

type eatableObject struct {
	Id uint64
	*erutan.Component_SpaceComponent
}

type EatableSystem struct {
	objects octree.Octree
}

func NewEatableSystem() *EatableSystem {
	return &EatableSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(),
		utils.Config.GroundSize))}
}

func (e *EatableSystem) Add(id uint64,
	space *erutan.Component_SpaceComponent) {
	eo := eatableObject{id, space}
	o := octree.NewObjectCube(eo, eo.Position.Get(0), eo.Position.Get(1), eo.Position.Get(2), 1)
	e.objects.Insert(*o)
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (e *EatableSystem) Remove(o octree.Object) {
	e.objects.Remove(o)
}

func (e *EatableSystem) Update(dt float64) {
	/*
		for _, entity := range e.entities {
		}
	*/
}

func (e *EatableSystem) NotifyCallback(event utils.Event) {
	switch u := event.Value.(type) {
	case ObjectsCollided:
		a := u.a.Data.(collisionObject)
		b := u.b.Data.(collisionObject)
		// If an animal collided with me
		if a.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			b.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			newSc := b.Component_SpaceComponent
			p := protometry.RandomCirclePoint(*protometry.NewVectorN(utils.Config.GroundSize/2,
				utils.Config.GroundSize/2),
				utils.Config.GroundSize/2)
			newSc.Position = &p
			ManagerInstance.Watch.Notify(utils.Event{Value: ObjectPhysicsUpdated{object: u.b, newSc: *newSc, dt: u.dt}})

		}

		if b.BehaviourType == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			a.BehaviourType == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			newSc := a.Component_SpaceComponent
			p := protometry.RandomCirclePoint(*protometry.NewVectorN(utils.Config.GroundSize/2,
				utils.Config.GroundSize/2),
				utils.Config.GroundSize/2)
			newSc.Position = &p
			ManagerInstance.Watch.Notify(utils.Event{Value: ObjectPhysicsUpdated{object: u.a, newSc: *newSc, dt: u.dt}})
		}

	}
}
