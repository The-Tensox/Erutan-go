package game

import (
	"github.com/The-Tensox/erutan/cfg"
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
	*erutan.Component_NetworkBehaviourComponent
}

func (a AnyObject) ID() uint64 {
	return a.Id
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
		cfg.Global.Logic.GroundSize*1000))}
}

func (e *EatableSystem) Add(id uint64,
	space *erutan.Component_SpaceComponent) {
	eo := eatableObject{id, space}
	o := octree.NewObjectCube(eo, eo.Position.Get(0), eo.Position.Get(1), eo.Position.Get(2), 1)
	if !e.objects.Insert(*o) {
		utils.DebugLogf("Failed to insert %v", o.ToString())
	}
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (e *EatableSystem) Remove(o octree.Object) {
	e.objects.Remove(o)
}

func (e *EatableSystem) Update(dt float64) {
}

func (e *EatableSystem) Handle(event utils.Event) {
	switch u := event.Value.(type) {
	case utils.ObjectsCollided:
		a := u.A.Data.(collisionObject)
		b := u.B.Data.(collisionObject)
		// If an animal collided with me
		if a.Tag == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			b.Tag == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			newSc := b.Component_SpaceComponent
			p := protometry.RandomCirclePoint(*protometry.NewVectorN(cfg.Global.Logic.GroundSize, cfg.Global.Logic.GroundSize),
				cfg.Global.Logic.GroundSize)
			newSc.Position = &p
			ManagerInstance.Watch.NotifyAll(utils.Event{Value: utils.ObjectPhysicsUpdated{Object: u.B, NewSc: *newSc, Dt: u.Dt}})
		}

		if b.Tag == erutan.Component_BehaviourTypeComponent_ANIMAL &&
			a.Tag == erutan.Component_BehaviourTypeComponent_VEGETATION {
			// Teleport somewhere else
			newSc := a.Component_SpaceComponent
			p := protometry.RandomCirclePoint(*protometry.NewVectorN(cfg.Global.Logic.GroundSize, cfg.Global.Logic.GroundSize),
				cfg.Global.Logic.GroundSize)
			newSc.Position = &p
			ManagerInstance.Watch.NotifyAll(utils.Event{Value: utils.ObjectPhysicsUpdated{Object: u.A, NewSc: *newSc, Dt: u.Dt}})
		}
	}
}
